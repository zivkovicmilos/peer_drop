import { Box, Typography } from '@material-ui/core';
import { makeStyles } from '@material-ui/core/styles';
import InfoRoundedIcon from '@material-ui/icons/InfoRounded';
import LockRoundedIcon from '@material-ui/icons/LockRounded';
import SecurityRoundedIcon from '@material-ui/icons/SecurityRounded';
import clsx from 'clsx';
import { FC } from 'react';
import { ReactComponent as ChecklistRoundedIcon } from '../../../shared/assets/icons/checklist_black_24dp.svg';
import { ReactComponent as BlackBar } from '../../../shared/assets/icons/line_vertical.svg';
import theme from '../../../theme/theme';
import { ENewWorkspaceStep } from '../../pages/NewWorkspace/newWorkspace.types';
import { INewWorkspaceStepsProps } from './newWorkspaceSteps.types';

const NewWorkspaceSteps: FC<INewWorkspaceStepsProps> = (props) => {
  const { currentStep, currentStepIndex, steps } = props;

  const classes = useStyles();

  const isActive = (localIndex: number) => {
    return currentStepIndex >= localIndex;
  };

  const isActiveBar = (localIndex: number) => {
    return currentStepIndex > localIndex;
  };

  const renderIcon = (step: ENewWorkspaceStep, localIndex: number) => {
    switch (step) {
      case ENewWorkspaceStep.PARAMS:
        return (
          <InfoRoundedIcon
            style={{
              width: '30px',
              height: 'auto',
              fill: isActive(localIndex)
                ? 'black'
                : theme.palette.custom.transparentBlack
            }}
          />
        );
      case ENewWorkspaceStep.SECURITY:
        return (
          <SecurityRoundedIcon
            style={{
              width: '30px',
              height: 'auto',
              fill: isActive(localIndex)
                ? 'black'
                : theme.palette.custom.transparentBlack
            }}
          />
        );
      case ENewWorkspaceStep.PERMISSIONS:
        return (
          <LockRoundedIcon
            style={{
              width: '30px',
              height: 'auto',
              fill: isActive(localIndex)
                ? 'black'
                : theme.palette.custom.transparentBlack
            }}
          />
        );
      case ENewWorkspaceStep.REVIEW:
        return (
          <ChecklistRoundedIcon
            style={{
              width: '30px',
              height: 'auto',
              fill: isActive(localIndex)
                ? 'black'
                : theme.palette.custom.transparentBlack
            }}
          />
        );
    }
  };

  const renderSection = (step: ENewWorkspaceStep, index: number) => {
    return (
      <Box
        key={`step-${index}`}
        display={'flex'}
        flexDirection={'column'}
        alignItems={'center'}
        mb={1}
      >
        <Box mb={0.5}>{renderIcon(step, index)}</Box>
        <Typography
          className={clsx(classes.stepTitle, {
            [classes.inactiveStep]: !isActive(index)
          })}
        >
          {step}
        </Typography>
      </Box>
    );
  };

  const renderStep = (step: ENewWorkspaceStep, index: number) => {
    if (index != steps.length - 1) {
      return (
        <Box display={'flex'} flexDirection={'column'} alignItems={'center'}>
          {renderSection(step, index)}
          <Box>
            <BlackBar
              style={{
                width: 'auto',
                height: '60px',
                fill: isActiveBar(index)
                  ? 'black'
                  : theme.palette.custom.transparentBlack
              }}
            />
          </Box>
        </Box>
      );
    } else {
      return renderSection(step, index);
    }
  };

  return (
    <Box display={'flex'} flexDirection={'column'}>
      {steps.map((step, index) => {
        {
          return renderStep(step, index);
        }
      })}
    </Box>
  );
};

const useStyles = makeStyles(() => {
  return {
    stepTitle: {
      fontSize: '0.875rem',
      fontWeight: 500
    },
    inactiveStep: {
      color: theme.palette.custom.transparentBlack
    }
  };
});

export default NewWorkspaceSteps;
