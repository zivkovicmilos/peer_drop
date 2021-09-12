import { Box, Typography } from '@material-ui/core';
import { makeStyles } from '@material-ui/core/styles';
import { FC, useContext } from 'react';
import NewWorkspaceContext from '../../../context/NewWorkspaceContext';
import { ReactComponent as SuccessImage } from '../../../shared/assets/icons/undraw_well_done_i2wr.svg';
import theme from '../../../theme/theme';
import MnemonicCopy from '../../atoms/MnemonicCopy/MnemonicCopy';
import NewWorkspaceInfo from '../NewWorkspaceInfo/NewWorkspaceInfo';
import { INewWorkspaceSuccessProps } from './newWorkspaceSuccess.types';

const NewWorkspaceSuccess: FC<INewWorkspaceSuccessProps> = () => {
  const {
    workspaceName,
    workspaceType,
    workspaceMnemonic,
    accessControlType,
    accessControl,
    permissions
  } = useContext(NewWorkspaceContext);

  const classes = useStyles();

  return (
    <Box display={'flex'} flexDirection={'column'} width={'100%'}>
      <Box
        display={'flex'}
        flexDirection={'column'}
        width={'100%'}
        alignItems={'center'}
      >
        <SuccessImage
          style={{
            width: '30%',
            height: 'auto'
          }}
        />
        <Box mt={4}>
          <Typography className={classes.workspaceName}>
            {workspaceName}
          </Typography>
        </Box>
        <MnemonicCopy mnemonic={workspaceMnemonic} />
      </Box>
      <Box display={'flex'} flexDirection={'column'} width={'50%'} mt={4}>
        <NewWorkspaceInfo
          accessControl={accessControl}
          accessControlType={accessControlType}
          permissions={permissions}
          workspaceName={workspaceName}
          workspaceType={workspaceType}
        />
      </Box>
    </Box>
  );
};

const useStyles = makeStyles(() => {
  return {
    workspaceName: {
      fontWeight: 700,
      fontSize: theme.typography.pxToRem(24)
    }
  };
});
 
export default NewWorkspaceSuccess;
