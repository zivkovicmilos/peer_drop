import { Box } from '@material-ui/core';
import { makeStyles } from '@material-ui/core/styles';
import { FC } from 'react';
import SingleWorkspace from '../../atoms/SingleWorkspace/SingleWorkspace';
import { IWorkspaceListProps } from './workspaceList.types';

const WorkspaceList: FC<IWorkspaceListProps> = () => {
  const classes = useStyles();

  const workspaceList = [
    {
      title: 'Polygon'
    },
    {
      title: 'ZP Projekat'
    },
    {
      title: 'KRIK'
    },
    {
      title: 'Al Jazeera'
    },
    {
      title: 'Work group'
    },
    {
      title: 'Example #1'
    },
    {
      title: 'Example #2'
    },
    {
      title: 'Example #3'
    }
  ];

  return (
    <Box className={classes.workspaceListWrapper}>
      {workspaceList.map((workspace, index) => {
        return (
          <SingleWorkspace
            key={`${workspace.title}-item-${index}`}
            title={workspace.title}
          />
        );
      })}
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
