export const getServiceStatus = (token, serviceName) => (
  getRequest(`/api/status?serviceName=${serviceName}`, token))

export const getServiceStatuses = token => (
  getRequest(`/api/statuses`, token))

// getRequest creates and executes a get request
const getRequest = (route, token) => (
  fetch(route, makeRequestObject('GET', token))
  .then(checkReponse)
  .then(res => res.json())
);

// makeRequestObject creates a request object with HTTP method and headers.
const makeRequestObject = (method, token) => ({
  method,
  headers: {
    'Accept': 'application/json',
    'Content-Type': 'application/json',
    'Authorization': `Basic ${token}`,
  },
})

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
