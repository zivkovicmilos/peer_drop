import { Box } from '@material-ui/core';
import { makeStyles } from '@material-ui/core/styles';
import { FC, useEffect } from 'react';
import NoData from '../../atoms/NoData/NoData';
import Pagination from '../../atoms/Pagination/Pagination';
import usePagination from '../../atoms/Pagination/pagination.hook';
import SingleWorkspace from '../../atoms/SingleWorkspace/SingleWorkspace';
import { IWorkspaceListProps } from './workspaceList.types';

const WorkspaceList: FC<IWorkspaceListProps> = () => {
  const classes = useStyles();

  const workspaceList = [
    {
      id: '1',
      title: 'Polygon'
    },
    {
      id: '2',
      title: 'ZP Projekat'
    },
    {
      id: '3',
      title: 'KRIK'
    },
    {
      id: '4',
      title: 'Al Jazeera'
    },
    {
      id: '5',
      title: 'Work group'
    },
    {
      id: '6',
      title: 'Example #1'
    },
    {
      id: '7',
      title: 'Example #2'
    },
    {
      id: '8',
      title: 'Example #3'
    }
  ];

  const { page, count, setCount, limit, handlePageChange } = usePagination({
    limit: 8
  });

  useEffect(() => {
    setCount(workspaceList.length);
  }, []);

  return (
    <Box className={classes.workspaceListWrapper}>
      {workspaceList.map((workspace, index) => {
        return (
          <SingleWorkspace
            key={`${workspace.title}-item-${index}`}
            title={workspace.title}
            id={workspace.id}
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

      {workspaceList.length < 1 && (
        <Box width={'100%'} mt={8}>
          <NoData text={'No workspaces found'} />{' '}
        </Box>
      )}
    </Box>
  );
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
