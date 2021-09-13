import { Box, Button, Chip, IconButton } from '@material-ui/core';
import { makeStyles } from '@material-ui/core/styles';
import ArrowBackRoundedIcon from '@material-ui/icons/ArrowBackRounded';
import MoreVertRoundedIcon from '@material-ui/icons/MoreVertRounded';
import { FC, useState } from 'react';
import { useParams } from 'react-router-dom';
import { ENewWorkspaceType } from '../../../context/newWorkspaceContext.types';
import { ReactComponent as UploadIcon } from '../../../shared/assets/icons/file_upload_black_24dp.svg';
import theme from '../../../theme/theme';
import Link from '../../atoms/Link/Link';
import PageTitle from '../../atoms/PageTitle/PageTitle';
import SectionTitle from '../../atoms/SectionTitle/SectionTitle';
import SuggestedList from '../../molecules/SuggestedList/SuggestedList';
import { IFileInfo } from '../../molecules/SuggestedList/suggestedList.types';
import ViewWorkspaceFiles from '../../molecules/ViewWorkspaceFiles/ViewWorkspaceFiles';
import { IViewWorkspaceProps, IWorkspaceInfo } from './viewWorkspace.types';

const ViewWorkspace: FC<IViewWorkspaceProps> = () => {
  const { workspaceId } = useParams() as { workspaceId: string };

  // TODO update this with the server
  const [workspaceInfo, setWorkspaceInfo] = useState<IWorkspaceInfo>({
    id: workspaceId,
    name: 'Polygon',
    type: ENewWorkspaceType.SEND_RECEIVE
  });

  const classes = useStyles();

  const [suggestedFiles, setSuggestedFiles] = useState<{ files: IFileInfo[] }>({
    files: [
      {
        id: '1',
        name: 'File 1',
        extension: 'pdf'
      },
      {
        id: '2',
        name: 'File 2',
        extension: 'docx'
      },
      {
        id: '3',
        name: 'File 3',
        extension: 'xls'
      },
      {
        id: '4',
        name: 'File 4',
        extension: 'zip'
      },
      {
        id: '5',
        name: 'File 5',
        extension: 'txt'
      }
    ]
  });

  return (
    <Box display={'flex'} width={'100%'} flexDirection={'column'}>
      <Box display={'flex'} alignItems={'center'} width={'80%'}>
        {
          // TODO add unsaved changes modal
        }
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
          <PageTitle title={workspaceInfo.name} />
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
          <Button variant={'text'} startIcon={<UploadIcon />}>
            Upload
          </Button>
          <Box ml={2}>
            <Chip
              size={'small'}
              label={workspaceInfo.type}
              className={classes.chip}
            />
          </Box>
        </Box>
      </Box>
      <Box display={'flex'} flexDirection={'column'} mt={4}>
        <Box mb={4}>
          <SectionTitle title={'Suggested'} />
        </Box>
        <Box>
          <SuggestedList files={suggestedFiles.files} />
        </Box>
      </Box>
      <Box display={'flex'} flexDirection={'column'} mt={4}>
        <Box mb={2}>
          <SectionTitle title={'Workspace files'} />
        </Box>
        <Box>
          <ViewWorkspaceFiles workspaceInfo={workspaceInfo} />
        </Box>
      </Box>
    </Box>
  );
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
