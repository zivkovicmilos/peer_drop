import React from 'react';

export interface IStepButtonProps {
  text: string;
  onClick?: () => void;
  disabled?: boolean;
  variant: string

  shouldSubmit?: boolean;
}