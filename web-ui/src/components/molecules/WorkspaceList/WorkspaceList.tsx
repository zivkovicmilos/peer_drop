import { Box } from '@material-ui/core';
import { makeStyles } from '@material-ui/core/styles';
import { FC, useEffect, useState } from 'react';
import WorkspacesService from '../../../services/workspaces/workspacesService';
import { IWorkspaceListResponse } from '../../../services/workspaces/workspacesService.types';
import theme from '../../../theme/theme';
import LoadingIndicator from '../../atoms/LoadingIndicator/LoadingIndicator';
import NoData from '../../atoms/NoData/NoData';
import Pagination from '../../atoms/Pagination/Pagination';
import usePagination from '../../atoms/Pagination/pagination.hook';
import SingleWorkspace from '../../atoms/SingleWorkspace/SingleWorkspace';
import useSnackbar from '../Snackbar/useSnackbar.hook';
import { IWorkspaceListProps } from './workspaceList.types';

const WorkspaceList: FC<IWorkspaceListProps> = () => {
  const classes = useStyles();

  const [workspaces, setWorkspaces] = useState<IWorkspaceListResponse>(null);
  const { page, count, setCount, limit, handlePageChange } = usePagination({
    limit: 8
  });

  const { openSnackbar } = useSnackbar();

  useEffect(() => {
    const fetchWorkspaces = async () => {
      return await WorkspacesService.getWorkspaces({ page, limit });
    };

    fetchWorkspaces()
      .then((response) => {
        setCount(response.count);
        setWorkspaces(response);
      })
      .catch((err) => {
        openSnackbar('Unable to fetch workspaces', 'error');
      });
  }, []);

  if (workspaces && count > 0) {
    return (
      <Box className={classes.workspaceListWrapper}>
        {workspaces.workspaceWrappers.map((workspace, index) => {
          return (
            <SingleWorkspace
              key={`${workspace.workspaceName}-item-${index}`}
              title={workspace.workspaceName}
              id={workspace.workspaceMnemonic}
              mnemonic={workspace.workspaceMnemonic}
            />
          );
        })}
        <Box
          display={'flex'}
          alignItems={'center'}
          justifyContent={'center'}
          width={'100%'}
          mt={8}
        >
          <Pagination
            count={count}
            limit={limit}
            page={page}
            onPageChange={handlePageChange}
          />
        </Box>

        {workspaces.workspaceWrappers.length < 1 && (
          <Box width={'100%'} mt={8}>
            <NoData text={'No workspaces found'} />{' '}
          </Box>
        )}
      </Box>
    );
  } else if (count == 0) {
    return <NoData text={'No workspaces found'} />;
  } else {
    return <LoadingIndicator style={{ color: theme.palette.primary.main }} />;
  }
};

const useStyles = makeStyles(() => {
  return {
    workspaceListWrapper: {
      display: 'flex',
      flexWrap: 'wrap',
      width: '100%',
      marginTop: '-50px',
      marginLeft: '-50px'
    }
  };
});

export default WorkspaceList;
