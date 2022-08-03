import { Modal as BootstrapModal } from 'react-bootstrap';
const Modal = (props: any) => {
  const { show, title, header, body, footer, onHide, ...otherProps } = props;

  return (
    show && (
      <BootstrapModal show={show} onHide={onHide} {...otherProps}>
        {(title || header) && (
          <BootstrapModal.Header closeButton>
            {header}
            {title && <BootstrapModal.Title>{title}</BootstrapModal.Title>}
          </BootstrapModal.Header>
        )}
        {body && <BootstrapModal.Body>{body}</BootstrapModal.Body>}
        {footer && <BootstrapModal.Footer>{footer}</BootstrapModal.Footer>}
      </BootstrapModal>
    )
  );
};

export default Modal;
