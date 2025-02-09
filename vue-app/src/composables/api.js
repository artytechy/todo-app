import router from "@/router";
import { useUserStore } from "@/stores/user.js";

export function useApi() {
  const user = useUserStore();

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
    if (response.status === 401) {
      user.logout();
      router.push("/login");
    }

    return response.json();
  }

  async function post(url, payload) {
    const response = await fetch(url, requestOptions("POST", payload));
    if (response.status === 401) {
      user.logout();
      router.push("/login");
    }
    
    return response.json();
  }

  return { get, post };
}
