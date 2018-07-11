export const getServiceStatus = (token, serviceName) => (
  getRequest(`/api/status/${serviceName}`, token))

export const getServiceStatuses = token => (
  getRequest(`/api/statuses`, token))

// postRequestJSON creates and executes a post request and parses JSON response
const postRequest = (route, token, body = {}) => (
  fetch(route, makeRequestObjectWithBody('POST', token, body))
  .then(checkReponse)
  .then(res => res.json())
);

// getRequest creates and executes a get request
const getRequest = (route, token) => (
  fetch(route, makeRequestObject('GET', token))
  .then(checkReponse)
  .then(res => res.json())
);

// makeRequestObjectWithBody returns a reqest object to be passed to fetch in
// order to make a request of the supplied method
const makeRequestObjectWithBody = (method, token, body) => ({
  ...makeRequestObject(method, token),
  body: JSON.stringify(body)
})

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
