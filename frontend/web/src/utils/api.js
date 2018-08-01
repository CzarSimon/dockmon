export const getServiceStatus = (username, password, serviceName) => (
  getRequest(`/api/status?serviceName=${serviceName}`, username, password));

export const getServiceStatuses = (username, password) => (
  getRequest(`/api/statuses`, username, password));

export const login = (username, password) => (
  postRequest('/api/login', username, password));

// getRequest creates and executes a GET request.
const getRequest = (route, username, password) => (
  fetch(route, makeRequestObject('GET', username, password))
    .then(checkReponse)
    .then(res => res.json())
);

// postRequest creates and executes a POST request.
const postRequest = (route, username, password) => (
  fetch(route, makeRequestObject('POST', username, password))
    .then(checkReponse)
    .then(res => res.json())
);

// makeRequestObject creates a request object with HTTP method and headers.
const makeRequestObject = (method, username, password) => {
  const token = btoa(`${username}:${password}`);
  return {
    method,
    headers: {
      'Accept': 'application/json',
      'Content-Type': 'application/json',
      'Authorization': `Basic ${token}`,
    },
  }
}

// checkReponse checks whether a fetch response was ok, throws an error if not
const checkReponse = response => {
  if (response.ok) {
    return response
  } else {
    let error = {};
    error.response = response;
    throw error;
  }
};
