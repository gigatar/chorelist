export const unixtTimeToLocal = unixTime => {
  if (typeof unixTime !== "number") {
    return "";
  }

  const date = new Date(unixTime * 1000);

  return date.toLocaleDateString() + " " + date.toLocaleTimeString();
};
