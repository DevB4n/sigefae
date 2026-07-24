package crear_correos

import (
	"strings"

	"sigefae/internal/correo"
	"sigefae/procesos_aplicacion/descarga_correos/graph"
)

func ProcesarYCrear(msg graph.Message, correoSvc *correo.Service) error {
	req := correo.CreateRequest{
		Asunto:         msg.Subject,
		De:             msg.From.EmailAddress.Address,
		Para:           joinRecipients(msg.ToRecipients),
		FechaRecepcion: msg.ReceivedDateTime,
		IDMensaje:      msg.ID,
		Cuerpo:         msg.Body.Content,
		Cc:             joinRecipients(msg.CcRecipients),
		Bcc:            joinRecipients(msg.BccRecipients),
		ReplyTo:        joinRecipients(msg.ReplyTo),
		IDEstado:       1,
	}

	_, err := correoSvc.Create(req)
	return err
}

func joinRecipients(recipients []graph.Recipient) string {
	if len(recipients) == 0 {
		return ""
	}
	var addresses []string
	for _, r := range recipients {
		if r.EmailAddress.Address != "" {
			addresses = append(addresses, r.EmailAddress.Address)
		}
	}
	return strings.Join(addresses, ", ")
}
