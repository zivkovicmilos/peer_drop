import { Button } from '@material-ui/core';
import { FC } from 'react';
import { IActionButtonProps } from './actionButton.types';

const ActionButton: FC<IActionButtonProps> = (props) => {
  const { text, onClick, startIcon, disabled = false } = props;

  return (
    <Button
      variant={'contained'}
      onClick={onClick}
      className={'actionButtonRounded'}
      type={'submit'}
      color={'primary'}
      startIcon={startIcon}
      disabled={disabled}
    >
      {text}
    </Button>
  );
};

export default ActionButton;
