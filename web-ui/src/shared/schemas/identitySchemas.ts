import * as yup from 'yup';

const generateKeyValidationSchema = yup.object({
  generateName: yup.string().defined('Name is required'),
  generateEmail: yup
    .string()
    .email('Email is not valid')
    .defined('Email is required')
});

export default generateKeyValidationSchema;
