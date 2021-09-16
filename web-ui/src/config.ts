export const { REACT_APP_CLIENT_API_BASE_URL } = process.env;

if (!REACT_APP_CLIENT_API_BASE_URL) {
  throw Error('REACT_APP_CLIENT_API_BASE_URL property not found in .env');
}

export default {
  API_BASE_URL: REACT_APP_CLIENT_API_BASE_URL
};
