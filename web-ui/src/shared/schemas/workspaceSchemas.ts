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

const digitsOnly = (value: any) => /^\d+$/.test(value);

const nwPermissionsSchema = yup.object({
  autoCloseWorkspace: yup.boolean(),
  autoCloseDate: yup.date().when('autoCloseWorkspace', {
    is: (autoCloseWorkspace: boolean) => {
      return autoCloseWorkspace;
    },
    then: yup.date().required('Auto close date is required'),
    otherwise: yup.date()
  }),

  // Peer limit cannot be enforced if the security
  // type is: Specific contacts
  enforcePeerLimit: yup.boolean(),
  peerLimit: yup.string().when('enforcePeerLimit', {
    is: (enforcePeerLimit: boolean) => {
      return enforcePeerLimit;
    },
    then: yup
      .string()
      .required('Number of peers is required')
      .test('Value test', 'Only numbers allowed', digitsOnly)
      .test('Value test 2', 'Only positive numbers allowed', function (value) {
        return +value > 0;
      })
      .test('Value test 3', 'Maximum number of peers is 100', function (value) {
        return +value <= 100;
      }),
    otherwise: yup.string()
  })
});

export { nwSecurityPasswordSchema, nwParametersSchema, nwPermissionsSchema };
