import os
from dotenv import load_dotenv
from fastapi import FastAPI, Query, Body
from pydantic import BaseModel, Field
from typing import List, Dict, Any, Optional
from datetime import date, datetime
from influxdb3 import InfluxConnectorV3

# --- Configuration ---
INFLUXDB_URL = os.getenv("INFLUXDB_URL")
INFLUXDB_TOKEN = os.getenv("INFLUXDB_TOKEN")
INFLUXDB_ORG = os.getenv("INFLUXDB_ORG")
INFLUXDB_DATABASE = os.getenv("INFLUXDB_DATABASE")

# --- InfluxDB Client ---
# Use a placeholder if the env vars are not set
if INFLUXDB_URL and INFLUXDB_TOKEN and INFLUXDB_ORG:
    # Initialize V3 Connector
    db = InfluxConnectorV3(INFLUX_HOST, INFLUX_TOKEN, INFLUX_ORG, INFLUX_DB)
else:
    db = None # No client if config is missing

# --- FastAPI App ---
app = FastAPI(
    title="VitalStream API",
    description="API for the VitalStream Health Dashboard, using FastAPI and InfluxDB.",
    version="1.0.0",
)

# --- Pydantic Models ---

class Metric(BaseModel):
    measurement: str
    tags: Dict[str, str]
    fields: Dict[str, Any]
    timestamp: str

class IngestData(BaseModel):
    metrics: List[Metric]

class SummaryResponse(BaseModel):
    steps: int
    distance: float
    activeCalories: float
    basalCalories: float
    dietaryCalories: float

class HeartRateResponse(BaseModel):
    time: str
    value: int

class BloodPressureResponse(BaseModel):
    time: str
    systolic: int
    diastolic: int
    category: str

class GlucoseResponse(BaseModel):
    time: str
    value: int

class SleepResponse(BaseModel):
    date: str
    totalDuration: float
    deepSleep: float
    remSleep: float
    lightSleep: float
    awake: float
    efficiency: float

class WorkoutResponse(BaseModel):
    id: str
    time: str
    name: str
    duration: int
    calories: int
    type: str
    avgHr: int

class DietaryTrendResponse(BaseModel):
    date: str
    calories: int
    protein: int
    carbs: int
    fat: int
    trend: int

class MealResponse(BaseModel):
    name: str
    desc: str
    cal: int

class BodyCompositionResponse(BaseModel):
    time: str
    weight: float
    bodyFat: float


# --- API Endpoints ---

@app.post("/api/v1/ingest", status_code=202)
async def ingest_data(data: IngestData = Body(...)):
    # In a real implementation, you would write this data to InfluxDB
    # For now, we just accept it.
    if db:
        # The influxdb3-python client is synchronous for writes,
        # but we can run it in a thread pool executor if needed in a truly async app.
        # For this shell, direct call is fine.
        # client.write(record=data.metrics) # This is a conceptual example
        pass
    return {"message": "Data ingestion started."}

@app.get("/api/v1/summary", response_model=SummaryResponse)
async def get_summary(date: Optional[date] = Query(None)):
    # SQL Query
    sql_query = """
    SELECT *
    FROM "daily_totals"
    WHERE time > now() - interval '1d'
    ORDER BY time ASC
    """

    today = datetime.now(ZoneInfo("America/New_York")).date()
    df = db.get_dataframe(sql_query).with_columns(pl.col('date').str.to_date('%Y-%m-%d')).filter(pl.col('date').eq(today))

    steps_query = df.filter((pl.col('metric')=='step_count') & (pl.col('source')=='RingConn')).select('value')
    distance_query = df.filter((pl.col('metric')=='walking_running_distance')).select('value')
    act_cal_query = df.filter((pl.col('metric')=='active_energy') & (pl.col('source')=='RingConn')).select('value')
    b_cal_query = df.filter((pl.col('metric')=='basal_energy_burned') & (pl.col('source')=='RingConn')).select('value')
    calories_query = df.filter((pl.col('metric')=='dietary_energy')).select('value')

    steps = 0.0 if steps_query.is_empty() else steps_query[0,0]
    distance = 0.0 if distance_query.is_empty() else distance_query[0,0]
    act_cal = 0.0 if act_cal_query.is_empty() else act_cal_query[0,0]
    b_cal = 0.0 if b_cal_query.is_empty() else b_cal_query[0,0]
    calories = 0.0 if calories_query.is_empty() else calories_query[0,0]

    # Make response
    response = {
        "steps": steps,
        "distance": distance,
        "activeCalories": act_cal,
        "basalCalories": b_cal,
        "dietaryCalories": calories
    }


    return response

@app.get("/api/v1/vitals/hr", response_model=List[HeartRateResponse])
async def get_heart_rate(date: Optional[date] = Query(None)):
    # SQL Query
    sql_query = """
    SELECT *
    FROM "heart_rate"
    WHERE time > now() - interval '1d'
    ORDER BY time ASC
    """

    today = datetime.now(ZoneInfo("America/New_York")).date()
    df = db.get_dataframe(sql_query)

    df = df.with_columns(
        pl.col("time")
        .dt.replace_time_zone("UTC")
        .dt.convert_time_zone("America/New_York")
        .dt.replace_time_zone(None)
    )

    df = df.group_by_dynamic("time", every="10m").agg(
        pl.col("avg").mean(),
        pl.col("max").max(),
        pl.col("min").min(),
    )

    result = df.select(['time', 'avg']).with_columns(pl.col('time').dt.strftime('%H:%M')).rename({'avg':'value'}).to_dicts()

    return result

@app.get("/api/v1/vitals/bp", response_model=List[BloodPressureResponse])
async def get_blood_pressure(end_date: Optional[date] = Query(None)):
    today = datetime.now(ZoneInfo("America/New_York"))
    start = today.isoformat()
    if not end_date:
        end = (today - timedelta(days=90)).isoformat()
    else:
        end = datetime.strptime(end_date, '%Y-%m-%d').isoformat()


    # SQL Query
    sql_query = f"""
    SELECT *
    FROM "blood_pressure"
    WHERE time > '{end}' and time <= '{start}'
    ORDER BY time ASC
    """

    df = db.get_dataframe(sql_query)
    df = df.with_columns(
        pl.when((pl.col("systolic") > 180) | (pl.col("diastolic") > 120))
        .then(pl.lit("Hypertensive Crisis"))
        .when((pl.col("systolic") >= 140) | (pl.col("diastolic") >= 90))
        .then(pl.lit("Hypertension Stage 2"))
        .when((pl.col("systolic") >= 130) | (pl.col("diastolic") >= 80))
        .then(pl.lit("Hypertension Stage 1"))
        .when((pl.col("systolic") >= 120) & (pl.col("diastolic") < 80))
        .then(pl.lit("Elevated"))
        .otherwise(pl.lit("Normal"))
        .cast(pl.Categorical)
        .alias("category")
    )
    result = df.with_columns(pl.col('time').dt.strftime('%h %e')).select(['time', 'systolic', 'diastolic', 'category']).to_dicts()

    return result

@app.get("/api/v1/vitals/glucose", response_model=List[GlucoseResponse])
async def get_glucose(end_date: Optional[date] = Query(None)):
    today = datetime.now(ZoneInfo("America/New_York"))
    start = today.isoformat()
    if not end_date:
        end = (today - timedelta(days=90)).isoformat()
    else:
        end = datetime.strptime(end_date, '%Y-%m-%d').isoformat()

    print (f'start = {start}; end = {end}')

    # SQL Query
    sql_query = f"""
    SELECT *
    FROM "blood_glucose"
    WHERE time > '{end}' and time <= '{start}'
    ORDER BY time ASC
    """

    df = db.get_dataframe(sql_query)
    result = df.with_columns(pl.col('time').dt.strftime('%h %e')).rename({'qty':'value'}).select(['time', 'value']).to_dicts()

    return result

@app.get("/api/v1/sleep", response_model=List[SleepResponse])
async def get_sleep(end_date: Optional[date] = Query(None)):
    today = datetime.now(ZoneInfo("America/New_York"))
    start = today.isoformat()
    if not end_date:
        end = (today - timedelta(days=90)).isoformat()
    else:
        end = datetime.strptime(end_date, '%Y-%m-%d').isoformat()


    # SQL Query
    sql_query = f"""
    SELECT *
    FROM "sleep_analysis"
    WHERE time > '{end}' and time <= '{start}'
    ORDER BY time ASC
    """

    df = db.get_dataframe(sql_query)
    result = df.rename({'time':'date', 'totalSleep':'totalDuration', 'deep':'deepSleep', 'rem':'remSleep', 'core':'lightSleep'}) \
    .select(['date', 'totalDuration', 'deepSleep', 'remSleep', 'lightSleep']).with_columns(
        pl.col('date').dt.strftime('%h %e'),
        pl.lit(95).alias('efficiency')
    ).to_dicts()

    return result

@app.get("/api/v1/workouts", response_model=List[WorkoutResponse])
async def get_workouts(date: Optional[date] = Query(None)):
    today = datetime.now(ZoneInfo("America/New_York"))
    start = today.isoformat()
    if not date:
        end = (today - timedelta(days=90)).isoformat()
    else:
        end = datetime.strptime(date, '%Y-%m-%d').isoformat()


    # SQL Query
    sql_query = f"""
    SELECT *
    FROM "workout"
    WHERE time > '{end}' and time <= '{start}'
    ORDER BY time ASC
    """

    df = db.get_dataframe(sql_query)
    df = df.select(['workout_id', 'time', 'workout_name', 'duration', 'active_energy_value']).with_columns(
        pl.col('workout_name').alias('type'),
        (pl.col('duration') / 60).cast(pl.Int64),
        pl.col('active_energy_value').cast(pl.Int64),
    )

    # Heart Rate
    # SQL Query
    sql_query = f"""
    SELECT *
    FROM "workout_heart_rate"
    WHERE time > '{end}' and time <= '{start}'
    ORDER BY time ASC
    """

    df2 = db.get_dataframe(sql_query)
    df2 = df2.group_by('workout_id').agg(pl.col('avg').mean().cast(pl.Int64))
    result = df.join(df2, on='workout_id', how='left').with_columns(
        pl.col('time').dt.strftime('%Y-%m-%d %H:%M')
    ).rename({'workout_name':'name', 'active_energy_value':'calories', 'avg':'avgHr'}).to_dicts()


    return result

@app.get("/api/v1/dietary/trends", response_model=List[DietaryTrendResponse])
async def get_dietary_trends(end_date: Optional[date] = Query(None)):
    today = datetime.now(ZoneInfo("America/New_York"))
    start = today.isoformat()
    if not end_date:
        end = (today - timedelta(days=90)).isoformat()
        tend = (today - timedelta(days=97)).isoformat()
    else:
        end = datetime.strptime(end_date, '%Y-%m-%d').isoformat()
        tend = (datetime.strptime(end_date, '%Y-%m-%d') - timedelta(days=7)).isoformat()

    #---------------------
    # Calories
    #---------------------

    # SQL Query
    sql_query = f"""
    SELECT *
    FROM "dietary_energy"
    WHERE time > '{end}' and time <= '{start}'
    ORDER BY time ASC
    """

    df = db.get_dataframe(sql_query)
    df = df.with_columns(
        pl.col('time').dt.convert_time_zone(time_zone="America/New_York")
    ).with_columns(
        pl.col('time').dt.strftime('%Y-%m-%d').alias('day')
    ).group_by('day').agg(pl.col('qty').sum().cast(pl.Int64).alias('calories')).sort('day')

    #---------------------
    # Protein
    #---------------------

    # SQL Query
    sql_query = f"""
    SELECT *
    FROM "protein"
    WHERE time > '{end}' and time <= '{start}'
    ORDER BY time ASC
    """

    df2 = db.get_dataframe(sql_query)
    df2 = df2.with_columns(
        pl.col('time').dt.convert_time_zone(time_zone="America/New_York")
    ).with_columns(
        pl.col('time').dt.strftime('%Y-%m-%d').alias('day')
    ).group_by('day').agg(pl.col('qty').sum().cast(pl.Int64).alias('protein')).sort('day')

    #---------------------
    # Carbs
    #---------------------

    # SQL Query
    sql_query = f"""
    SELECT *
    FROM "carbohydrates"
    WHERE time > '{end}' and time <= '{start}'
    ORDER BY time ASC
    """

    df3 = db.get_dataframe(sql_query)
    df3 = df3.with_columns(
        pl.col('time').dt.convert_time_zone(time_zone="America/New_York")
    ).with_columns(
        pl.col('time').dt.strftime('%Y-%m-%d').alias('day')
    ).group_by('day').agg(pl.col('qty').sum().cast(pl.Int64).alias('carbs')).sort('day')

    #---------------------
    # Fat
    #---------------------

    # SQL Query
    sql_query = f"""
    SELECT *
    FROM "total_fat"
    WHERE time > '{end}' and time <= '{start}'
    ORDER BY time ASC
    """

    df4 = db.get_dataframe(sql_query)
    df4 = df4.with_columns(
        pl.col('time').dt.convert_time_zone(time_zone="America/New_York")
    ).with_columns(
        pl.col('time').dt.strftime('%Y-%m-%d').alias('day')
    ).group_by('day').agg(pl.col('qty').sum().cast(pl.Int64).alias('fat')).sort('day')

    #---------------------
    # Trend
    #---------------------

    # SQL Query
    sql_query = f"""
    SELECT *
    FROM "dietary_energy"
    WHERE time > '{tend}' and time <= '{start}'
    ORDER BY time ASC
    """

    df5 = db.get_dataframe(sql_query)
    df5 = df5.with_columns(
        pl.col('time').dt.convert_time_zone(time_zone="America/New_York")
    ).with_columns(
        pl.col('time').dt.strftime('%Y-%m-%d').alias('day')
    ).group_by('day').agg(pl.col('qty').sum().alias('trend')).with_columns(
        pl.col('trend').rolling_mean(window_size=7, min_samples=3).cast(pl.Int64)
    ).sort('day')

    ### Combine all ####
    df = df.join(df2, on='day', how='left').join(df3, on='day', how='left').join(df4, on='day', how='left').join(df5, on='day', how='left')
    result = df.with_columns(
        pl.col('day').str.to_datetime('%Y-%m-%d').dt.strftime('%h %e'),
        pl.col('trend').fill_null(strategy='forward')
    ).rename({'day':'time'}).to_dicts()

    return result

@app.get("/api/v1/dietary/meals/today", response_model=List[MealResponse])
async def get_meals_today(date: Optional[date] = Query(None)):
    return [
        { "name": "Breakfast", "desc": "Oatmeal, Berries, Whey", "cal": 420 },
        { "name": "Lunch", "desc": "Chicken Salad, Quinoa", "cal": 580 }
    ]

@app.get("/api/v1/body/composition", response_model=List[BodyCompositionResponse])
async def get_body_composition(end_date: Optional[date] = Query(None)):
    today = datetime.now(ZoneInfo("America/New_York"))
    start = today.isoformat()
    if not end_date:
        end = (today - timedelta(days=90)).isoformat()
    else:
        end = datetime.strptime(end_date, '%Y-%m-%d').isoformat()

    #---------------------
    # Weight
    #---------------------
    # SQL Query
    sql_query = f"""
    SELECT *
    FROM "weight_body_mass"
    WHERE time > '{end}' and time <= '{start}'
    ORDER BY time ASC
    """

    df = db.get_dataframe(sql_query)
    df = df.rename({'qty':'weight'})

    #---------------------
    # BF%
    #---------------------

    # SQL Query
    sql_query = f"""
    SELECT *
    FROM "body_fat_percentage"
    WHERE time > '{end}' and time <= '{start}'
    ORDER BY time ASC
    """

    df2 = db.get_dataframe(sql_query)
    df2 = df2.rename({'qty':'bodyFat'})

    df = df.join(df2, on='time', how='left')
    result = df.drop_nulls().select(['time', 'weight', 'bodyFat']).with_columns(
        pl.col('time').dt.strftime('%h %e')
    ).to_dicts()

    return result
