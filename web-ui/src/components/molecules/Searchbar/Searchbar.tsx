import { InputBase, Paper } from '@material-ui/core';
import { makeStyles, Theme } from '@material-ui/core/styles';
import SearchRoundedIcon from '@material-ui/icons/SearchRounded';
import { FC, useContext } from 'react';
import SessionContext from '../../../context/SessionContext';
import theme from '../../../theme/theme';
import { ISearchBarProps } from './searchbar.types';

const Searchbar: FC<ISearchBarProps> = () => {
  const classes = useStyles();
  const { searchContext } = useContext(SessionContext);

  return (
    <Paper component={'form'} className={classes.searchBar}>
      <SearchRoundedIcon
        style={{
          fill: theme.palette.custom.transparentBlack
        }}
      />

      <InputBase className={classes.input} placeholder={'Search'} />
    </Paper>
  );
};

const useStyles = makeStyles((theme: Theme) => {
  return {
    searchBar: {
      display: 'flex',
      minWidth: '75%', // TODO change this to be more dynamic
      alignItems: 'center',
      borderRadius: '10px',
      padding: '10px 15px',
      backgroundColor: theme.palette.custom.darkGray,
      boxShadow: 'none',
      outline: 'none',
      '&:focus-within': {
        boxShadow: theme.palette.boxShadows.darker
      }
    },
    input: {
      marginLeft: theme.spacing(1),
      flex: 1
    }
  };
});

export default Searchbar;
