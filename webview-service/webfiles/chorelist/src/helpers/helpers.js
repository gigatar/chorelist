export const unixtTimeToLocal = unixTime => {
  if (typeof unixTime !== "number") {
    return "";
  }

  const date = new Date(unixTime * 1000);

  return date.toLocaleDateString() + " " + date.toLocaleTimeString();
};

export const validEmail = email => {
  let re = /^(([^<>()[\]\\.,;:\s@"]+(\.[^<>()[\]\\.,;:\s@"]+)*)|(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/;
  return re.test(email);
};
