import {
  Box,
  Button,
  Chip,
  ClickAwayListener,
  Grow,
  IconButton,
  MenuItem,
  MenuList,
  Paper,
  Popper,
  Tooltip
} from '@material-ui/core';
import { makeStyles } from '@material-ui/core/styles';
import ArrowBackRoundedIcon from '@material-ui/icons/ArrowBackRounded';
import ExitToAppRoundedIcon from '@material-ui/icons/ExitToAppRounded';
import MoreVertRoundedIcon from '@material-ui/icons/MoreVertRounded';
import React, {
  FC,
  Fragment,
  useContext,
  useEffect,
  useRef,
  useState
} from 'react';
import { useHistory, useParams } from 'react-router-dom';
import { ENewWorkspaceType } from '../../../context/newWorkspaceContext.types';
import SessionContext from '../../../context/SessionContext';
import WorkspacesService from '../../../services/workspaces/workspacesService';
import { IWorkspaceDetailedResponse } from '../../../services/workspaces/workspacesService.types';
import { ReactComponent as Clipboard } from '../../../shared/assets/icons/content_paste_black_24dp.svg';
import { ReactComponent as UploadIcon } from '../../../shared/assets/icons/file_upload_black_24dp.svg';
import { ReactComponent as Loupe } from '../../../shared/assets/icons/loupe.svg';
import CommonUtils from '../../../shared/utils/CommonUtils';
import theme from '../../../theme/theme';
import Link from '../../atoms/Link/Link';
import LoadingIndicator from '../../atoms/LoadingIndicator/LoadingIndicator';
import NoData from '../../atoms/NoData/NoData';
import PageTitle from '../../atoms/PageTitle/PageTitle';
import SectionTitle from '../../atoms/SectionTitle/SectionTitle';
import useSnackbar from '../../molecules/Snackbar/useSnackbar.hook';
import SuggestedList from '../../molecules/SuggestedList/SuggestedList';
import ViewWorkspaceFiles from '../../molecules/ViewWorkspaceFiles/ViewWorkspaceFiles';
import { IViewWorkspaceProps } from './viewWorkspace.types';

const ViewWorkspace: FC<IViewWorkspaceProps> = () => {
  const { workspaceMnemonic } = useParams() as { workspaceMnemonic: string };

  const { userIdentity } = useContext(SessionContext);

  const [isOwner, setIsOwner] = useState<boolean>(false);

  const [workspaceDetailed, setWorkspaceDetailed] =
    useState<IWorkspaceDetailedResponse>(null);

  const { openSnackbar } = useSnackbar();

  const history = useHistory();

  const setIfOwner = (publicKeyIDs: string[]) => {
    for (let i = 0; i < publicKeyIDs.length; i++) {
      if (publicKeyIDs[i] == userIdentity.publicKeyID) {
        setIsOwner(true);
      }
    }

    setIsOwner(false);
  };

  const fetchWorkspaceDetailed = () => {
    const fetchDetailed = async () => {
      return await WorkspacesService.getWorkspaceFiles(
        CommonUtils.formatMnemonic(workspaceMnemonic)
      );
    };

    fetchDetailed()
      .then((response) => {
        setWorkspaceDetailed(response);

        setIfOwner(response.workspaceOwnerKeyIDs);
      })
      .catch((err) => {
        openSnackbar('Unable to fetch workspace files', 'error');

        history.push('/workspaces');
      });
  };

  // Set the updater
  useEffect(() => {
    fetchWorkspaceDetailed();

    const workspaceInfoUpdater = setInterval(() => {
      fetchWorkspaceDetailed();
    }, 5000);

    return () => clearInterval(workspaceInfoUpdater);
  }, []);

  const classes = useStyles();

  const convertWorkspaceType = (type: string) => {
    switch (type) {
      case 'send-only':
        return ENewWorkspaceType.SEND_ONLY;
      case 'receive-only':
        return ENewWorkspaceType.RECEIVE_ONLY;
      default:
        return ENewWorkspaceType.SEND_RECEIVE;
    }
  };

  const renderSuggestedFiles = () => {
    if (workspaceDetailed && workspaceDetailed.workspaceFiles.length > 0) {
      return (
        <SuggestedList
          files={workspaceDetailed.workspaceFiles.slice(0, 5)}
          workspaceMnemonic={workspaceMnemonic}
        />
      );
    }

    return (
      <NoData
        text={'No files to suggest'}
        icon={
          <Loupe
            style={{
              width: '50px',
              height: 'auto'
            }}
          />
        }
      />
    );
  };

  const imageRef: any = useRef();

  const showOpenFileDialog = () => {
    imageRef.current.click();
  };

  const uploadFile = async (formData: FormData) => {
    return await WorkspacesService.uploadWorkspaceFile(formData);
  };

  const handleUpload = (event: any) => {
    const fileObject = event.target.files[0];
    if (!fileObject) return;

    const formData = new FormData();

    formData.append('workspaceFile', fileObject, fileObject.name);
    formData.set('mnemonic', workspaceMnemonic);

    uploadFile(formData)
      .then((response) => {
        openSnackbar('File successfully shared to the workspace', 'success');
      })
      .catch((err) => {
        openSnackbar('Failed to upload file', 'error');
      });
  };

  const renderWorkspaceActions = () => {
    let workspaceType = convertWorkspaceType(workspaceDetailed.workspaceType);

    if (isOwner) {
      if (
        workspaceType == ENewWorkspaceType.SEND_ONLY ||
        workspaceType == ENewWorkspaceType.SEND_RECEIVE
      ) {
        return (
          <Fragment>
            <input
              ref={imageRef}
              type="file"
              style={{ display: 'none' }}
              accept="*"
              onChange={handleUpload}
            />
            <Button
              variant={'text'}
              startIcon={<UploadIcon />}
              onClick={showOpenFileDialog}
            >
              Upload
            </Button>
          </Fragment>
        );
      } else {
        return (
          <Tooltip title={'Workspace is receive-only'}>
            <Button disabled variant={'text'} startIcon={<UploadIcon />}>
              Upload
            </Button>
          </Tooltip>
        );
      }
    } else {
      if (workspaceType == ENewWorkspaceType.SEND_ONLY) {
        return (
          <Tooltip title={'Workspace is send-only'}>
            <Button disabled variant={'text'} startIcon={<UploadIcon />}>
              Upload
            </Button>
          </Tooltip>
        );
      } else {
        return (
          <Fragment>
            <input
              ref={imageRef}
              type="file"
              style={{ display: 'none' }}
              accept="*"
              onChange={handleUpload}
            />
            <Button
              variant={'text'}
              startIcon={<UploadIcon />}
              onClick={showOpenFileDialog}
              disabled={
                workspaceType == ENewWorkspaceType.RECEIVE_ONLY && isOwner
              }
            >
              Upload
            </Button>
          </Fragment>
        );
      }
    }
  };

  const handleWorkspaceCopy = () => {
    navigator.clipboard.writeText(
      CommonUtils.unformatMnemonic(workspaceMnemonic)
    );

    openSnackbar('Copied to clipboard!', 'success');
  };

  const handleWorkspaceLeave = () => {
    const leaveWorkspace = async () => {
      return await WorkspacesService.leaveWorkspace(workspaceMnemonic);
    };

    leaveWorkspace()
      .then((response) => {
        history.push('/workspaces');

        openSnackbar('Workspace left', 'success');
      })
      .catch((err) => {
        openSnackbar('Unable to leave workspace', 'error');
      });
  };

  const [open, setOpen] = useState(false);
  const anchorRef = useRef<HTMLDivElement>(null);

  const handleToggle = () => {
    setOpen((prevOpen) => !prevOpen);
  };

  const handleClose = (event: any) => {
    if (
      anchorRef.current &&
      anchorRef.current.contains(event.target as HTMLElement)
    ) {
      return;
    }

    setOpen(false);
  };

  if (workspaceDetailed) {
    return (
      <Box display={'flex'} width={'100%'} flexDirection={'column'}>
        <Box display={'flex'} alignItems={'center'} width={'80%'}>
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
            <PageTitle title={workspaceDetailed.workspaceName} />
          </Box>
          <div
            style={{
              marginLeft: 1
            }}
            ref={anchorRef}
          >
            <IconButton
              classes={{
                root: 'iconButtonRoot'
              }}
            >
              <MoreVertRoundedIcon
                style={{
                  fill: 'black',
                  height: '18px',
                  width: 'auto'
                }}
                onClick={handleToggle}
              />
            </IconButton>
            <Popper
              placement={'bottom-start'}
              open={open}
              anchorEl={anchorRef.current}
              role={undefined}
              transition
            >
              {({ TransitionProps, placement }) => (
                <Grow
                  {...TransitionProps}
                  style={{
                    transformOrigin:
                      placement === 'bottom' ? 'center top' : 'center bottom'
                  }}
                >
                  <Paper>
                    <ClickAwayListener onClickAway={handleClose}>
                      <MenuList
                        style={{
                          background: 'white',
                          borderRadius: '5px',
                          boxShadow: 'none'
                        }}
                      >
                        <MenuItem
                          className={classes.userMenuItem}
                          onClick={() => {
                            handleWorkspaceCopy();
                            handleToggle();
                          }}
                        >
                          <Box display={'flex'}>
                            <Clipboard />
                            <Box ml={1}>Share</Box>
                          </Box>
                        </MenuItem>
                        <MenuItem
                          className={classes.userMenuItem}
                          onClick={() => {
                            handleWorkspaceLeave();
                            handleToggle();
                          }}
                        >
                          <Box display={'flex'}>
                            <ExitToAppRoundedIcon />
                            <Box ml={1}>Leave</Box>
                          </Box>
                        </MenuItem>
                      </MenuList>
                    </ClickAwayListener>
                  </Paper>
                </Grow>
              )}
            </Popper>
          </div>

          <Box ml={'auto'} display={'flex'} alignItems={'center'}>
            {renderWorkspaceActions()}
            <Box ml={2}>
              <Chip
                size={'small'}
                label={convertWorkspaceType(workspaceDetailed.workspaceType)}
                className={classes.chip}
              />
            </Box>
          </Box>
        </Box>
        <Box display={'flex'} flexDirection={'column'} mt={4}>
          <Box mb={4}>
            <SectionTitle title={'Suggested'} />
          </Box>
          <Box>{renderSuggestedFiles()}</Box>
        </Box>
        <Box display={'flex'} flexDirection={'column'} mt={4}>
          <Box mb={2}>
            <SectionTitle title={'Workspace files'} />
          </Box>
          <Box>
            <ViewWorkspaceFiles workspaceInfo={workspaceDetailed} />
          </Box>
        </Box>
      </Box>
    );
  } else {
    return <LoadingIndicator style={{ color: theme.palette.primary.main }} />;
  }
};

const useStyles = makeStyles(() => {
  return {
    chip: {
      border: `1px solid #A6B0B9`,
      borderRadius: '4px',
      backgroundColor: '#EFF2F5',
      color: '#A6B0B9',
      fontSize: theme.typography.pxToRem(14),
      fontWeight: 500
    },
    userMenuItem: {
      fontFamily: 'Montserrat',
      fontWeight: 500,
      fontSize: '0.875rem',
      color: 'black'
    }
  };
});

export default ViewWorkspace;
