import * as yup from 'yup';
import { ENWAccessControl } from '../../context/newWorkspaceContext.types';

const nwParametersSchema = yup.object({
  workspaceName: yup.string().defined('Workspace name is required')
});

const nwSecurityPasswordSchema = yup.object({
  accessControlType: yup.string().defined(),
  password: yup.string().when('accessControlType', {
    is: (accessControlType: string) => {
      return accessControlType == ENWAccessControl.PASSWORD;
    },
    then: yup.string().required('Password is required'),
    otherwise: yup.string()
  }),
  passwordConfirm: yup.string().when('accessControlType', {
    is: (accessControlType: string) => {
      return accessControlType == ENWAccessControl.PASSWORD;
    },
    then: yup
      .string()
      .required('Confirm password is required')
      .test('Password match test', 'Passwords must match', function (value) {
        return value === this.parent.password;
      }),
    otherwise: yup.string()
  })
});

export { nwSecurityPasswordSchema, nwParametersSchema };
