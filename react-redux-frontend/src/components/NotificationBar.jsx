import { Snackbar, Alert } from '@mui/material';
import { useAppDispatch, useAppSelector } from '../hooks/useRedux';
import { removeNotification } from '../store/slices/uiSlice';

function NotificationBar() {
  const dispatch = useAppDispatch();
  const notifications = useAppSelector((state) => state.ui.notifications);

  const handleClose = (id) => {
    dispatch(removeNotification(id));
  };

  return (
    <>
      {notifications.map((notification, index) => (
        <Snackbar
          key={notification.id}
          open={true}
          autoHideDuration={6000}
          onClose={() => handleClose(notification.id)}
          anchorOrigin={{ vertical: 'top', horizontal: 'right' }}
          style={{ top: `${80 + index * 70}px` }}
        >
          <Alert
            onClose={() => handleClose(notification.id)}
            severity={notification.severity || 'info'}
            sx={{ width: '100%' }}
          >
            {notification.message}
          </Alert>
        </Snackbar>
      ))}
    </>
  );
}

export default NotificationBar;

