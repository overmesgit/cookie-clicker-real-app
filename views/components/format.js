export function formatNumber(num) {
    const suffixes = [
        {value: 1e12, symbol: "T"},
        {value: 1e9, symbol: "B"},
        {value: 1e6, symbol: "M"},
        // {value: 1e3, symbol: "K"}
    ];

    const suffix = suffixes.find(s => Math.abs(num) >= s.value);

    if (!suffix) {
        return num;
    }

    const formattedValue = parseFloat((num / suffix.value).toFixed(6));
    return `${formattedValue}${suffix.symbol}`;
}