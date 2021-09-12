import { Box, IconButton, TextField, Typography } from '@material-ui/core';
import { makeStyles } from '@material-ui/core/styles';
import ArrowBackRoundedIcon from '@material-ui/icons/ArrowBackRounded';
import RefreshRoundedIcon from '@material-ui/icons/RefreshRounded';
import WifiTetheringRoundedIcon from '@material-ui/icons/WifiTetheringRounded';
import { useFormik } from 'formik';
import React, { FC, useEffect, useState } from 'react';
import { useHistory } from 'react-router-dom';
import {
  ENewWorkspaceType,
  ENWAccessControl
} from '../../../context/newWorkspaceContext.types';
import { joinWorkspaceSchema } from '../../../shared/schemas/workspaceSchemas';
import theme from '../../../theme/theme';
import ActionButton from '../../atoms/ActionButton/ActionButton';
import FormTitle from '../../atoms/FormTitle/FormTitle';
import Link from '../../atoms/Link/Link';
import LoadingIndicator from '../../atoms/LoadingIndicator/LoadingIndicator';
import PageTitle from '../../atoms/PageTitle/PageTitle';
import JoinWorkspaceModal from '../../molecules/JoinWorkspaceModal/JoinWorkspaceModal';
import NewWorkspaceInfo from '../../molecules/NewWorkspaceInfo/NewWorkspaceInfo';
import { IJoinWorkspaceProps, WorkspaceInfo } from './joinWorkspace.types';

const JoinWorkspace: FC<IJoinWorkspaceProps> = () => {
  const blankWorkspaceInfo: WorkspaceInfo = {
    workspaceName: 'Dummy workspace',
    workspaceType: ENewWorkspaceType.SEND_ONLY,
    accessControl: { password: '123' },
    permissions: {
      autocloseWorkspace: {
        active: false
      },

      enforcePeerLimit: {
        active: false
      },

      additionalOwners: {
        active: false
      }
    },
    accessControlType: ENWAccessControl.PASSWORD
  };

  const connectionFormik = useFormik({
    initialValues: {
      workspaceMnemonic: ''
    },
    validationSchema: joinWorkspaceSchema,
    onSubmit: (values, { resetForm }) => {
      // TODO trigger connection fetch

      setJoinState(JOIN_STATE.FETCH_INFO);

      setTimeout(() => {
        setWorkspaceInfo(blankWorkspaceInfo);
      }, 5000);
    }
  });

  enum JOIN_STATE {
    ENTER_MNEMONIC,
    FETCH_INFO,
    SHOW_INFO
  }

  const [buttonText, setButtonText] = useState<string>('Connect');
  const [buttonIcon, setButtonIcon] = useState<React.ReactNode>(
    <WifiTetheringRoundedIcon />
  );

  const [workspaceInfo, setWorkspaceInfo] =
    useState<WorkspaceInfo>(blankWorkspaceInfo);
  const [joinState, setJoinState] = useState<JOIN_STATE>(
    JOIN_STATE.ENTER_MNEMONIC
  );

  const classes = useStyles();

  useEffect(() => {
    if (joinState == JOIN_STATE.FETCH_INFO) {
      setJoinState(JOIN_STATE.SHOW_INFO);
    }
  }, [workspaceInfo]);

  useEffect(() => {
    if (joinState == JOIN_STATE.SHOW_INFO) {
      setButtonText('Refresh');
      setButtonIcon(<RefreshRoundedIcon />);
    }
  }, [joinState]);

  const handleWorkspaceJoin = () => {
    setConfirmOpen(true);
  };

  const renderLowerSection = () => {
    switch (joinState) {
      case JOIN_STATE.ENTER_MNEMONIC:
        return (
          <Typography className={classes.subText}>
            Please enter a workspace mnemonic
          </Typography>
        );
      case JOIN_STATE.FETCH_INFO:
        return (
          <Box display={'flex'} alignItems={'center'}>
            <Box mr={2}>
              <LoadingIndicator size={32} />
            </Box>
            <Typography className={classes.subText}>
              Fetching workspace information
            </Typography>
          </Box>
        );
      case JOIN_STATE.SHOW_INFO:
        if (workspaceInfo != null) {
          return (
            <Box width={'60%'} display={'flex'} flexDirection={'column'}>
              <NewWorkspaceInfo
                workspaceName={workspaceInfo.workspaceName}
                workspaceType={workspaceInfo.workspaceType}
                accessControl={workspaceInfo.accessControl}
                permissions={workspaceInfo.permissions}
                accessControlType={workspaceInfo.accessControlType}
              />
              <Box mt={2}>
                <ActionButton
                  text={'Join workspace'}
                  shouldSubmit={false}
                  startIcon={<WifiTetheringRoundedIcon />}
                  onClick={handleWorkspaceJoin}
                />
              </Box>
            </Box>
          );
        } else {
          return <Box>Dummy</Box>;
        }
    }
  };

  const [confirmOpen, setConfirmOpen] = useState<boolean>(false);

  const history = useHistory();

  const handleConfirm = (confirmed: boolean) => {
    setConfirmOpen(false);

    history.push('/workspaces');
  };

  return (
    <Box display={'flex'} flexDirection={'column'} width={'100%'}>
      <Box display={'flex'} alignItems={'center'} width={'100%'}>
        <Link to={'/workspaces'}>
          <IconButton
            classes={{
              root: 'iconButtonRoot'
            }}
          >
            <ArrowBackRoundedIcon
              style={{
                fill: 'black'
              }}
            />
          </IconButton>
        </Link>
        <Box ml={2}>
          <PageTitle title={'Join Workspace'} />
        </Box>
      </Box>
      <Box display={'flex'} width={'100%'} mt={4} flexDirection={'column'}>
        <form autoComplete={'off'} onSubmit={connectionFormik.handleSubmit}>
          <Box mb={2}>
            <FormTitle title={'Connection info'} />
          </Box>
          <Box display={'flex'} width={'100%'} alignItems={'center'}>
            <Box minHeight={'80px'} width={'30%'} display={'flex'}>
              <TextField
                id={'workspaceMnemonic'}
                label={'Workspace Mnemonic'}
                variant={'outlined'}
                fullWidth
                value={connectionFormik.values.workspaceMnemonic}
                onChange={connectionFormik.handleChange}
                onBlur={connectionFormik.handleBlur}
                error={
                  connectionFormik.touched.workspaceMnemonic &&
                  Boolean(connectionFormik.errors.workspaceMnemonic)
                }
                helperText={
                  connectionFormik.touched.workspaceMnemonic &&
                  connectionFormik.errors.workspaceMnemonic
                }
              />
            </Box>
            <Box display={'flex'} ml={4} mb={'24px'}>
              <ActionButton
                disabled={joinState == JOIN_STATE.FETCH_INFO}
                text={buttonText}
                startIcon={buttonIcon}
              />
            </Box>
          </Box>
        </form>
        <Box width={'100%'} display={'flex'} flexDirection={'column'}>
          <Box mb={2}>
            <FormTitle title={'Workspace info'} />
          </Box>
          <Box width={'100%'}>{renderLowerSection()}</Box>
        </Box>
      </Box>
      <JoinWorkspaceModal
        workspaceInfo={workspaceInfo}
        open={confirmOpen}
        handleConfirm={handleConfirm}
      />
    </Box>
  );
};

const useStyles = makeStyles(() => {
  return {
    subText: {
      fontSize: theme.typography.pxToRem(14)
    }
  };
});

export default JoinWorkspace;
