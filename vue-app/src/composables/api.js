import router from "@/router";
import { useUserStore } from "@/stores/user.js";
import { useCookies } from "./cookies.js";

export function useApi() {
  const user = useUserStore();
  const cookies = useCookies();

  function requestOptions(method = "GET", payload = null) {
    const headers = new Headers();
    headers.append("Content-Type", "application/json");
    headers.append("Authorization", "Bearer " + user.token);

    const options = {
      method,
      headers,
    };

    if (payload) {
      options.body = JSON.stringify(payload);
    }

    return options;
  }

  async function get(url) {
    const response = await fetch(url, requestOptions("GET"));
    checkResponse(response);
    return response.json();
  }

  async function post(url, payload) {
    const response = await fetch(url, requestOptions("POST", payload));
    checkResponse(response);
    return response.json();
  }

  function checkResponse(response) {
    if (response.status === 401) {
      user.logout();
      cookies.deleteCookie();
      router.push("/login");
    }
  }

  return { get, post };
}
