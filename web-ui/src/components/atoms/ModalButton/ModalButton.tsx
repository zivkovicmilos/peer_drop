import { Button } from '@material-ui/core';
import { FC } from 'react';
import { EModalButtonType, IModalButtonProps } from './modalButton.types';

const ModalButton: FC<IModalButtonProps> = (props) => {
  const { handleConfirm, type, text, margins } = props;

  const isConfirm = type == EModalButtonType.FILLED;

  return (
    <Button
      variant={isConfirm ? 'contained' : 'outlined'}
      color={'primary'}
      onClick={() => handleConfirm(isConfirm)}
      style={{
        margin: margins ? margins : 0,
        borderRadius: '15px'
      }}
    >
      {text}
    </Button>
  );
};

export default ModalButton;
