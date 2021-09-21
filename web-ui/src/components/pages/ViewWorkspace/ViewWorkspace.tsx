import { Box, Button, Chip, IconButton, Tooltip } from '@material-ui/core';
import { makeStyles } from '@material-ui/core/styles';
import ArrowBackRoundedIcon from '@material-ui/icons/ArrowBackRounded';
import MoreVertRoundedIcon from '@material-ui/icons/MoreVertRounded';
import { FC, Fragment, useEffect, useRef, useState } from 'react';
import { useHistory, useParams } from 'react-router-dom';
import { ENewWorkspaceType } from '../../../context/newWorkspaceContext.types';
import WorkspacesService from '../../../services/workspaces/workspacesService';
import { IWorkspaceDetailedResponse } from '../../../services/workspaces/workspacesService.types';
import { ReactComponent as UploadIcon } from '../../../shared/assets/icons/file_upload_black_24dp.svg';
import { ReactComponent as Loupe } from '../../../shared/assets/icons/loupe.svg';
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

  const [workspaceDetailed, setWorkspaceDetailed] =
    useState<IWorkspaceDetailedResponse>(null);

  const { openSnackbar } = useSnackbar();

  const history = useHistory();

  const fetchWorkspaceDetailed = () => {
    const fetchDetailed = async () => {
      return await WorkspacesService.getWorkspaceFiles(workspaceMnemonic);
    };

    fetchDetailed()
      .then((response) => {
        setWorkspaceDetailed(response);
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
        <SuggestedList files={workspaceDetailed.workspaceFiles.slice(0, 5)} />
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
          <Box ml={1}>
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
              />
            </IconButton>
          </Box>

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
    }
  };
});

export default ViewWorkspace;
