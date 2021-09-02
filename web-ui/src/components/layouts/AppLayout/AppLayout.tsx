import { Box, Drawer, SwipeableDrawer, useMediaQuery } from '@material-ui/core';
import { makeStyles, Theme } from '@material-ui/core/styles';
import { FC, useState } from 'react';
import theme from '../../../theme/theme';
import NavigationBar from '../../organisms/NavigationBar/NavigationBar';
import { IAppLayoutProps } from './appLayout.types';

const AppLayout: FC<IAppLayoutProps> = (props) => {
  const classes = useStyles();

  const [mobileOpen, setMobileOpen] = useState<boolean>(false);

  const handleDrawerToggle =
    (state: boolean) => (event: React.KeyboardEvent | React.MouseEvent) => {
      if (
        event &&
        event.type === 'keydown' &&
        ((event as React.KeyboardEvent).key === 'Tab' ||
          (event as React.KeyboardEvent).key === 'Shift')
      ) {
        return;
      }

      setMobileOpen(state);
    };

  const isBrowser = typeof window !== 'undefined';
  const iOS = isBrowser && /iPad|iPhone|iPod/.test(navigator.userAgent);

  const isMobile = useMediaQuery(theme.breakpoints.down('sm'));

  const renderPermanentDrawer = () => {
    return (
      <Drawer
        className={classes.drawer}
        variant="permanent"
        classes={{
          paper: classes.drawerPaper
        }}
        anchor="left"
      >
        <NavigationBar />
      </Drawer>
    );
  };

  const renderSwipeableDrawer = () => {
    return (
      <SwipeableDrawer
        anchor={'left'}
        open={mobileOpen}
        classes={{
          paper: classes.drawerPaper
        }}
        onClose={handleDrawerToggle(false)}
        onOpen={handleDrawerToggle(true)}
        disableBackdropTransition={!iOS}
        disableDiscovery={iOS}
      >
        <NavigationBar />
      </SwipeableDrawer>
    );
  };

  return (
    <Box display={'flex'} width={'100%'} flexGrow={1}>
      {isMobile ? renderSwipeableDrawer() : renderPermanentDrawer()}
      <Box
        py={{ xs: 0, md: 4 }}
        px={{ xs: 0, md: 8 }}
        width={'100%'}
        overflow={'auto'}
        className={classes.shadowWrapper}
      >
        {props.children}
      </Box>
    </Box>
  );
};

const drawerWidth = 300;

const useStyles = makeStyles((theme: Theme) => {
  return {
    shadowWrapper: {
      background: theme.palette.custom.white
    },
    drawer: {
      [theme.breakpoints.up('sm')]: {
        width: drawerWidth,
        flexShrink: 0
      }
    },
    drawerPaper: {
      width: drawerWidth,
      borderRight: '4px solid #F4F4F4'
    }
  };
});

export default AppLayout;
