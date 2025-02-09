export function useCookies() {
  function getCookie(name) {
    return document.cookie.split("; ").reduce((r, v) => {
      const parts = v.split("=");
      return parts[0] === name ? decodeURIComponent(parts[1]) : r;
    }, "");
  }

  function deleteCookie() {
    // Delete the cookie by setting an expired date
    document.cookie = `_SITE_DATA=; Path=/; SameSite=Strict; Secure; Expires=Thu, 01 Jan 1970 00:00:01 GMT;`;
  }

  function storeCookieWithUserData(data) {
    // Set expiration time for cookie (1 day)
    let date = new Date();
    date.setTime(date.getTime() + 24 * 60 * 60 * 1000);
    const expires = "expires=" + date.toUTCString();

    // Store user data securely in a cookie
    document.cookie = `_SITE_DATA=${JSON.stringify(
      data
    )}; ${expires}; path=/; SameSite=Strict; Secure;`;
  }

  return { getCookie, deleteCookie, storeCookieWithUserData };
}
