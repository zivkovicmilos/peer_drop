import { Button } from '@material-ui/core';
import { FC } from 'react';
import { IStepButtonProps } from './stepButton.types';

const StepButton: FC<IStepButtonProps> = (props) => {
  const {
    variant,
    text,
    onClick,
    disabled = false,
    shouldSubmit = true
  } = props;

  return (
    <Button
      variant={variant == 'contained' ? 'contained' : 'outlined'}
      onClick={onClick}
      type={shouldSubmit ? 'submit' : 'button'}
      color={'primary'}
      disabled={disabled}
    >
      {text}
    </Button>
  );
};

export default StepButton;
