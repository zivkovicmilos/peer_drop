export interface IPasswordStrengthProps {
  formik: any;
}

export enum EPasswordStrength {
  NULL,
  TOO_WEAK,
  WEAK,
  MEDIUM,
  STRONG
}
