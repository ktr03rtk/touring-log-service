type AlertProperties = {
  alertType: string;
  alertMessage: string;
};

const Alert = (props: AlertProperties) => {
  return (
    <div className={`alert ${props.alertType}`} role='alert'>
      {props.alertMessage}
    </div>
  );
};

export default Alert;
