import polars as pl
from influxdb_client_3 import InfluxDBClient3

class InfluxConnectorV3:
    def __init__(self, host, token, org, database):
        """
        Connects to InfluxDB 3.0 (IOx) utilizing Polars for
        zero-copy data handling.
        """
        self.client = InfluxDBClient3(
            host=host,
            token=token,
            org=org,
            database=database
        )

    def get_dataframe(self, sql_query):
        """
        Executes a SQL query and returns a Polars DataFrame.
        """
        try:
            # 1. Fetch data as a PyArrow Table (Native InfluxV3 format)
            arrow_table = self.client.query(query=sql_query, language="sql")

            # 2. Convert Arrow to Polars
            # This is extremely fast (zero-copy) compared to Pandas conversion
            df = pl.from_arrow(arrow_table)

            # Handle cases where result might be None or empty
            if df is None or df.is_empty():
                return pl.DataFrame()

            # 3. Sort by time
            # Polars does not have an "Index" like Pandas.
            # We explicitly sort by time to ensure lines draw correctly.
            if "time" in df.columns:
                df = df.sort("time")

            return df

        except Exception as e:
            # Catch connectivity issues or bad SQL syntax
            print(f"Query failed: {e}")
            return pl.DataFrame()

    def close(self):
        self.client.close()
