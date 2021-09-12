import { Box } from '@material-ui/core';
import { FC } from 'react';
import PageTitle from '../../atoms/PageTitle/PageTitle';
import WorkspaceList from '../../molecules/WorkspaceList/WorkspaceList';
import { IWorkspacesProps } from './workspaces.types';

const Workspaces: FC<IWorkspacesProps> = () => {
  return (
    <Box
      display={'flex'}
      flexDirection={'column'}
      width={'100%'}
      height={'100%'}
    >
      <Box display={'flex'} width={'100%'} mb={4}>
        <PageTitle title={'Workspaces'} />
      </Box>
      <WorkspaceList />
    </Box>
  );
};

export default Workspaces;
