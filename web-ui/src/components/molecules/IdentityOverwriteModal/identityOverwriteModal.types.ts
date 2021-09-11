export interface IIdentityOverwriteModalProps {
  publicKeyID: string;
  open: boolean;

  handleConfirm: (success: boolean) => void;
}
