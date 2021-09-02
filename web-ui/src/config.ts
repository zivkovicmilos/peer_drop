export const {
  CLIENT_API_BASE_URL,
} = process.env;

if (!CLIENT_API_BASE_URL) {
  throw Error('CLIENT_API_BASE_URL property not found in .env');
}



export default {
  API_BASE_URL: CLIENT_API_BASE_URL
};
