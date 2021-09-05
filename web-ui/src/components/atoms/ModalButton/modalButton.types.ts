export interface IModalButtonProps {
  handleConfirm: (success: boolean) => void;
  type: EModalButtonType;
  text: string;

  margins?: string;
}

export enum EModalButtonType {
  FILLED,
  OUTLINED
}