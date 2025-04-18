package client

func (cl *Client) Logout() error {
	stateXML, err := cl.GenerateStateXML()
	if err != nil {
		return err
	}

	_, err = cl.PostXMLWithSpecCtx(cl.TermURL, stateXML)
	if err != nil {
		return err
	}
	return nil
}
