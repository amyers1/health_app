export const formatNumber = (val: number | string) => {
    const num = typeof val === 'string' ? parseFloat(val) : val;

    if (isNaN(num)) {
        return val;
    }

    if (Number.isInteger(num)) {
        return num.toString();
    }

    return num.toFixed(2);
};
