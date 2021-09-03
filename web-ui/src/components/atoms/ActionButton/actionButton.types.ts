import React from 'react';

export interface IActionButtonProps {
  text: string;
  onClick?: () => void;
  startIcon?: React.ReactNode;
  disabled?: boolean;
}
