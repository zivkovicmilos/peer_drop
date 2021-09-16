import { Box, Button, IconButton, Menu, MenuItem } from '@material-ui/core';
import { makeStyles } from '@material-ui/core/styles';
import KeyboardArrowDownRoundedIcon from '@material-ui/icons/KeyboardArrowDownRounded';
import KeyboardArrowUpRoundedIcon from '@material-ui/icons/KeyboardArrowUpRounded';
import { FC, Fragment, useState } from 'react';
import {
  EIdentitySortDirection,
  EIdentitySortParam,
  IIdentitySortProps
} from './identitySort.types';

const IdentitySort: FC<IIdentitySortProps> = (props) => {
  const { activeSort, setActiveSort, sortDirection, setSortDirection } = props;
  const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null);

  const classes = useStyles();

  const handleClick = (event: any) => {
    setAnchorEl(event.currentTarget);
  };

  const handleClose = () => {
    setAnchorEl(null);
  };

  const renderSortDirection = () => {
    if (sortDirection == EIdentitySortDirection.ASC) {
      return <KeyboardArrowUpRoundedIcon className={classes.directionStyle} />;
    } else {
      return (
        <KeyboardArrowDownRoundedIcon className={classes.directionStyle} />
      );
    }
  };

  const handleDirectionChangeToggle = () => {
    setSortDirection(
      sortDirection == EIdentitySortDirection.ASC
        ? EIdentitySortDirection.DESC
        : EIdentitySortDirection.ASC
    );
  };

  const handleSortItemSelect = (type: EIdentitySortParam) => {
    setActiveSort(type);
  };

  const menuItems: EIdentitySortParam[] = [
    EIdentitySortParam.NAME,
    EIdentitySortParam.NUMBER_OF_WORKSPACES,
    EIdentitySortParam.CREATION_DATE,
    EIdentitySortParam.PUBLIC_KEY
  ];

  return (
    <Fragment>
      <Box display={'flex'} alignItems={'center'}>
        <Button onClick={handleClick}>{activeSort}</Button>
        <Box>
          <IconButton
            classes={{
              root: 'iconButtonRoot'
            }}
            onClick={() => handleDirectionChangeToggle()}
          >
            {renderSortDirection()}
          </IconButton>
        </Box>
      </Box>
      <Menu anchorEl={anchorEl} open={Boolean(anchorEl)} onClose={handleClose}>
        {menuItems.map((menuItem, index) => {
          return (
            <MenuItem
              key={`menuItem-${index}`}
              onClick={() => {
                handleSortItemSelect(menuItem);
                handleClose();
              }}
              className={classes.sortMenuItem}
            >
              {menuItem}
            </MenuItem>
          );
        })}
      </Menu>
    </Fragment>
  );
};

const useStyles = makeStyles(() => {
  return {
    directionStyle: {
      width: '20px',
      height: 'auto',
      fill: 'black'
    },
    sortMenuItem: {
      fontFamily: 'Montserrat',
      fontWeight: 500,
      fontSize: '0.875rem',
      color: 'black'
    }
  };
});

export default IdentitySort;