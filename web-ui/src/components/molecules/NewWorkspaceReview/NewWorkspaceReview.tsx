import { Box } from '@material-ui/core';
import { FC, useContext } from 'react';
import NewWorkspaceContext from '../../../context/NewWorkspaceContext';
import StepButton from '../../atoms/StepButton/StepButton';
import NewWorkspaceInfo from '../NewWorkspaceInfo/NewWorkspaceInfo';
import { INewWorkspaceReviewProps } from './newWorkspaceReview.types';

const NewWorkspaceReview: FC<INewWorkspaceReviewProps> = () => {
  const {
    handleNext,
    handleBack,
    step,
    setSectionTitle,
    setWorkspaceMnemonic
  } = useContext(NewWorkspaceContext);

  const handleConfirm = () => {
    // TODO alert the server
    setWorkspaceMnemonic('habit taste push wreck horse cotton');

    setSectionTitle('Workspace created!');
    handleNext();
  };

  return (
    <Box display={'flex'} flexDirection={'column'} width={'50%'}>
      <NewWorkspaceInfo />
      <Box display={'flex'} alignItems={'center'}>
        <Box mr={2}>
          <StepButton
            text={'Back'}
            variant={'outlined'}
            disabled={step < 1}
            shouldSubmit={false}
            onClick={() => {
              handleBack();
            }}
          />
        </Box>
        <StepButton
          text={'Confirm'}
          variant={'contained'}
          shouldSubmit={false}
          onClick={() => {
            handleConfirm();
          }}
        />
      </Box>
    </Box>
  );
};

export default NewWorkspaceReview;
