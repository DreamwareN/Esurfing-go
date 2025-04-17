package client

func (cl *Client) logout() error {
	stateXML, err := cl.GenerateStateXML()
	if err != nil {
		return err
	}

	_, err = cl.PostXMLWithoutCtx(cl.TermURL, stateXML)
	if err != nil {
		return err
	}
	return nil
}
