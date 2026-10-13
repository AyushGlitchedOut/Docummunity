//Other time formatting packages are pretty heavy and I dont need them so I made a simple one myself usign built in Intl library

const formatter = new Intl.RelativeTimeFormat("en", {
  numeric: "auto",
});

export default function TimeAgo(time: number): string {
  const seconds = Math.floor(Date.now() / 1000) - time;

  if (seconds < 60) return "Just Now";
  if (seconds < 3600)
    return formatter.format(-Math.floor(seconds / 60), "minutes");
  if (seconds < 86400)
    return formatter.format(-Math.floor(seconds / 3600), "hours");
  if (seconds < 604800)
    return formatter.format(-Math.floor(seconds / 86400), "days");
  if (seconds < 2592000)
    return formatter.format(-Math.floor(seconds / 604800), "weeks");
  if (seconds < 31536000)
    return formatter.format(-Math.floor(seconds / 2592000), "months");

  return formatter.format(-Math.floor(seconds / 31104000), "years");
}
