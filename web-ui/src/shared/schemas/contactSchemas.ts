import * as yup from 'yup';

const newContactValidationSchema = yup.object({
  name: yup.string().defined('Name is required')
});

export default newContactValidationSchema;
