export interface IJoinWorkspacePasswordModalProps {
  open: boolean;

  handleConfirm: (password: string, confirm: boolean) => void;
}