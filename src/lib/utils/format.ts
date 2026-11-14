/**
 * Format a `YYYY-MM-DD` date string as "Nov 9, 2026".
 * Splits manually instead of using Date parsing so the rendered
 * output is identical on the server and the client.
 */
export function formatDate(d: string): string {
  if (!d) return "";
  const parts = d.split("-");
  if (parts.length < 2) return parts[0];
  const months = [
    "Jan",
    "Feb",
    "Mar",
    "Apr",
    "May",
    "Jun",
    "Jul",
    "Aug",
    "Sep",
    "Oct",
    "Nov",
    "Dec",
  ];
  const monthIndex = parseInt(parts[1], 10) - 1;
  const month = months[monthIndex] || "";
  if (!month) return parts[0];
  if (parts.length >= 3 && parts[2]) {
    const day = parseInt(parts[2], 10);
    return isNaN(day) ? `${month} ${parts[0]}` : `${month} ${day}, ${parts[0]}`;
  }
  return `${month} ${parts[0]}`;
}
