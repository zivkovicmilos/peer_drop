import { Box, IconButton, TextField, Typography } from '@material-ui/core';
import { makeStyles } from '@material-ui/core/styles';
import ArrowBackRoundedIcon from '@material-ui/icons/ArrowBackRounded';
import RefreshRoundedIcon from '@material-ui/icons/RefreshRounded';
import WifiTetheringRoundedIcon from '@material-ui/icons/WifiTetheringRounded';
import { useFormik } from 'formik';
import React, { FC, useContext, useEffect, useState } from 'react';
import { useHistory } from 'react-router-dom';
import {
  ENewWorkspaceType,
  ENWAccessControl
} from '../../../context/newWorkspaceContext.types';
import SessionContext from '../../../context/SessionContext';
import WorkspacesService from '../../../services/workspaces/workspacesService';
import { IWorkspaceInfoResponse } from '../../../services/workspaces/workspacesService.types';
import { joinWorkspaceSchema } from '../../../shared/schemas/workspaceSchemas';
import CommonUtils from '../../../shared/utils/CommonUtils';
import theme from '../../../theme/theme';
import ActionButton from '../../atoms/ActionButton/ActionButton';
import FormTitle from '../../atoms/FormTitle/FormTitle';
import Link from '../../atoms/Link/Link';
import LoadingIndicator from '../../atoms/LoadingIndicator/LoadingIndicator';
import PageTitle from '../../atoms/PageTitle/PageTitle';
import JoinWorkspaceModal from '../../molecules/JoinWorkspaceModal/JoinWorkspaceModal';
import JoinWorkspacePasswordModal from '../../molecules/JoinWorkspacePasswordModal/JoinWorkspacePasswordModal';
import NewWorkspaceInfo from '../../molecules/NewWorkspaceInfo/NewWorkspaceInfo';
import useSnackbar from '../../molecules/Snackbar/useSnackbar.hook';
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

  const { openSnackbar } = useSnackbar();

  const connectionFormik = useFormik({
    initialValues: {
      workspaceMnemonic: ''
    },
    validationSchema: joinWorkspaceSchema,
    onSubmit: (values, { resetForm }) => {
      const fetchWorkspaceInfo = async () => {
        return await WorkspacesService.getWorkspaceInfo(
          CommonUtils.formatMnemonic(values.workspaceMnemonic)
        );
      };

      setJoinState(JOIN_STATE.FETCH_INFO);

      fetchWorkspaceInfo()
        .then((response) => {
          setWorkspaceInfo(response);
        })
        .catch((err) => {
          openSnackbar('Unable to fetch workspace info', 'error');
        });
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
    useState<IWorkspaceInfoResponse>(null);
  const [joinState, setJoinState] = useState<JOIN_STATE>(
    JOIN_STATE.ENTER_MNEMONIC
  );

  const classes = useStyles();

  const convertWorkspaceType = (type: string): ENewWorkspaceType => {
    switch (type) {
      case 'send-only':
        return ENewWorkspaceType.SEND_ONLY;
      case 'receive-only':
        return ENewWorkspaceType.RECEIVE_ONLY;
      default:
        return ENewWorkspaceType.SEND_RECEIVE;
    }
  };

  const convertAccessControlType = (type: string): ENWAccessControl => {
    if (type == 'password') {
      return ENWAccessControl.PASSWORD;
    } else {
      return ENWAccessControl.SPECIFIC_CONTACTS;
    }
  };

  const convertAccessControl = (
    workspaceInfo: IWorkspaceInfoResponse
  ): { contacts: string[] } | { password: string } => {
    if (workspaceInfo.passwordHash) {
      return {
        password: workspaceInfo.passwordHash
      };
    } else {
      return {
        contacts: workspaceInfo.contactsWrapper.contactPublicKeys
      };
    }
  };

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

  const { userIdentity } = useContext(SessionContext);

  const handleWorkspaceJoinInner = (password?: string) => {
    const joinWorkspace = async () => {
      return await WorkspacesService.joinWorkspace({
        mnemonic: workspaceInfo.mnemonic,
        password: password ? password : '',
        publicKeyID: userIdentity.publicKeyID
      });
    };

    joinWorkspace()
      .then((response) => {
        openSnackbar('Successfully joined the workspace', 'success');
        setConfirmOpen(true);

        history.push('/workspaces');
      })
      .catch((err) => {
        openSnackbar('Unable to join the workspace', 'error');
      });
  };

  const [passwordModalOpen, setPasswordModalOpen] = useState<boolean>(false);

  const handlePasswordConfirm = (password: string, confirm: boolean) => {
    setPasswordModalOpen(false);

    if (confirm) {
      if (password == '' || !password) {
        openSnackbar('Invalid password', 'error');

        return;
      }

      handleWorkspaceJoinInner(password);
    }
  };

  const handleWorkspaceJoin = () => {
    if (
      convertAccessControlType(workspaceInfo.securityType) ==
      ENWAccessControl.PASSWORD
    ) {
      // Fire up the password input modal
      setPasswordModalOpen(true);
    } else {
      handleWorkspaceJoinInner();
    }
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
                workspaceName={workspaceInfo.name}
                workspaceType={convertWorkspaceType(
                  workspaceInfo.workspaceType
                )}
                accessControl={convertAccessControl(workspaceInfo)}
                // permissions={workspaceInfo.permissions}
                accessControlType={convertAccessControlType(
                  workspaceInfo.securityType
                )}
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
      <JoinWorkspacePasswordModal
        open={passwordModalOpen}
        handleConfirm={handlePasswordConfirm}
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
