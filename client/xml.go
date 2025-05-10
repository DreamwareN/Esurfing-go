package client

import (
	"context"
	"encoding/xml"
	"io"
	"time"
)

const (
	UserAgentAndroid = "CCTP/android64_vpn/2093"
)

type TicketRequest struct {
	XMLName   xml.Name `xml:"request"`
	Text      string   `xml:",chardata"`
	UserAgent string   `xml:"user-agent"`
	ClientID  string   `xml:"client-id"`
	LocalTime string   `xml:"local-time"`
	HostName  string   `xml:"host-name"`
	Ipv4      string   `xml:"ipv4"`
	Ipv6      string   `xml:"ipv6"`
	Mac       string   `xml:"mac"`
	Ostag     string   `xml:"ostag"`
	Gwip      string   `xml:"gwip"`
}

type TicketResponse struct {
	XMLName xml.Name `xml:"response"`
	Text    string   `xml:",chardata"`
	Ticket  string   `xml:"ticket"`
	Expire  string   `xml:"expire"`
}

type LoginRequest struct {
	XMLName   xml.Name `xml:"request"`
	Text      string   `xml:",chardata"`
	UserAgent string   `xml:"user-agent"`
	ClientID  string   `xml:"client-id"`
	Ticket    string   `xml:"ticket"`
	LocalTime string   `xml:"local-time"`
	Userid    string   `xml:"userid"`
	Passwd    string   `xml:"passwd"`
}

type LoginResponse struct {
	XMLName    xml.Name `xml:"response"`
	Text       string   `xml:",chardata"`
	Userid     string   `xml:"userid"`
	KeepRetry  string   `xml:"keep-retry"`
	KeepURL    string   `xml:"keep-url"`
	TermURL    string   `xml:"term-url"`
	UserConfig struct {
		Text            string `xml:",chardata"`
		AgainstInterval string `xml:"against-interval"`
	} `xml:"user-config"`
	DomainConfig string `xml:"domain-config"`
}

type State struct {
	XMLName   xml.Name `xml:"request"`
	Text      string   `xml:",chardata"`
	UserAgent string   `xml:"user-agent"`
	ClientID  string   `xml:"client-id"`
	LocalTime string   `xml:"local-time"`
	HostName  string   `xml:"host-name"`
	Ipv4      string   `xml:"ipv4"`
	Ticket    string   `xml:"ticket"`
	Ipv6      string   `xml:"ipv6"`
	Mac       string   `xml:"mac"`
	Ostag     string   `xml:"ostag"`
}

type StateResponse struct {
	XMLName  xml.Name `xml:"response"`
	Text     string   `xml:",chardata"`
	Interval string   `xml:"interval"`
	Level    string   `xml:"level"`
}

type EConfig struct {
	XMLName       xml.Name `xml:"config"`
	Text          string   `xml:",chardata"`
	TicketURL     string   `xml:"ticket-url"`
	QueryURL      string   `xml:"query-url"`
	AuthURL       string   `xml:"auth-url"`
	StateURL      string   `xml:"state-url"`
	StateInterval string   `xml:"state-interval"`
	Auth          struct {
		Text    string `xml:",chardata"`
		Type    string `xml:"type"`
		AuthURL string `xml:"auth-url"`
		Default string `xml:"default"`
	} `xml:"auth"`
	Notify struct {
		Text     string `xml:",chardata"`
		Register string `xml:"register"`
	} `xml:"notify"`
	Against string `xml:"against"`
	Funcfg  struct {
		Text     string `xml:",chardata"`
		Province struct {
			Text string `xml:",chardata"`
			Name string `xml:"name,attr"`
		} `xml:"Province"`
		ManagerInternet struct {
			Text   string `xml:",chardata"`
			Enable string `xml:"enable,attr"`
		} `xml:"ManagerInternet"`
		QueryAnnouncement struct {
			Text   string `xml:",chardata"`
			Enable string `xml:"enable,attr"`
			URL    string `xml:"url,attr"`
		} `xml:"QueryAnnouncement"`
		QueryVerificateCodeStatus struct {
			Text   string `xml:",chardata"`
			URL    string `xml:"url,attr"`
			Enable string `xml:"enable,attr"`
		} `xml:"QueryVerificateCodeStatus"`
		QueryAuthCode struct {
			Text   string `xml:",chardata"`
			Enable string `xml:"enable,attr"`
			URL    string `xml:"url,attr"`
		} `xml:"QueryAuthCode"`
		PhoneMarketingData struct {
			Text   string `xml:",chardata"`
			Enable string `xml:"enable,attr"`
			URL    string `xml:"url,attr"`
		} `xml:"PhoneMarketingData"`
		SubErrorData struct {
			Text   string `xml:",chardata"`
			Enable string `xml:"enable,attr"`
			URL    string `xml:"url,attr"`
		} `xml:"SubErrorData"`
		SubErrorDataURL struct {
			Text   string `xml:",chardata"`
			Enable string `xml:"enable,attr"`
			URL    string `xml:"url,attr"`
		} `xml:"SubErrorDataURL"`
		ExperienceAccount struct {
			Text   string `xml:",chardata"`
			Enable string `xml:"enable,attr"`
			URL    string `xml:"url,attr"`
		} `xml:"ExperienceAccount"`
		Recharge struct {
			Text   string `xml:",chardata"`
			Enable string `xml:"enable,attr"`
			URL    string `xml:"url,attr"`
		} `xml:"Recharge"`
		PackageQuery struct {
			Text   string `xml:",chardata"`
			Enable string `xml:"enable,attr"`
			URL    string `xml:"url,attr"`
		} `xml:"PackageQuery"`
		UsedTime struct {
			Text   string `xml:",chardata"`
			Enable string `xml:"enable,attr"`
			URL    string `xml:"url,attr"`
		} `xml:"UsedTime"`
		NewSplashAndDialogAD struct {
			Text   string `xml:",chardata"`
			Enable string `xml:"enable,attr"`
			URL    string `xml:"url,attr"`
		} `xml:"NewSplashAndDialogAD"`
		PreAdvert struct {
			Text   string `xml:",chardata"`
			Type   string `xml:"type,attr"`
			Enable string `xml:"enable,attr"`
			URL    string `xml:"url,attr"`
		} `xml:"PreAdvert"`
		Advert struct {
			Text   string `xml:",chardata"`
			Enable string `xml:"enable,attr"`
			Type   string `xml:"type,attr"`
			URL    string `xml:"url,attr"`
		} `xml:"Advert"`
		ModifyPassword struct {
			Text   string `xml:",chardata"`
			Enable string `xml:"enable,attr"`
			URL    string `xml:"url,attr"`
		} `xml:"ModifyPassword"`
		ForgotPassword struct {
			Text   string `xml:",chardata"`
			Enable string `xml:"enable,attr"`
			URL    string `xml:"url,attr"`
		} `xml:"ForgotPassword"`
		PreventBrokenNetwork struct {
			Text   string `xml:",chardata"`
			Enable string `xml:"enable,attr"`
			URL    string `xml:"url,attr"`
		} `xml:"PreventBrokenNetwork"`
		FeedbackAndAdvice struct {
			Text   string `xml:",chardata"`
			Enable string `xml:"enable,attr"`
			URL    string `xml:"url,attr"`
		} `xml:"FeedbackAndAdvice"`
		MyFeedback struct {
			Text   string `xml:",chardata"`
			Enable string `xml:"enable,attr"`
			URL    string `xml:"url,attr"`
		} `xml:"MyFeedback"`
		UpdateCheck struct {
			Text   string `xml:",chardata"`
			Enable string `xml:"enable,attr"`
			URL    string `xml:"url,attr"`
		} `xml:"UpdateCheck"`
		Detect struct {
			Text   string `xml:",chardata"`
			Enable string `xml:"enable,attr"`
			URL    string `xml:"url,attr"`
			Weight string `xml:"weight,attr"`
		} `xml:"Detect"`
		ServiceHotline struct {
			Text     string `xml:",chardata"`
			Enable   string `xml:"enable,attr"`
			Phone    string `xml:"phone,attr"`
			WhiteSch string `xml:"whiteSch,attr"`
		} `xml:"ServiceHotline"`
		Vpn struct {
			Text   string `xml:",chardata"`
			Enable string `xml:"enable,attr"`
		} `xml:"Vpn"`
		Campus struct {
			Text   string `xml:",chardata"`
			Enable string `xml:"enable,attr"`
			URL    string `xml:"url,attr"`
		} `xml:"Campus"`
		CampusInterface struct {
			Text   string `xml:",chardata"`
			Enable string `xml:"enable,attr"`
			URL    string `xml:"url,attr"`
		} `xml:"CampusInterface"`
		CampusNewsTime struct {
			Text string `xml:",chardata"`
			Time string `xml:"time,attr"`
		} `xml:"CampusNewsTime"`
		SafetyDataCollectionURL struct {
			Text   string `xml:",chardata"`
			Enable string `xml:"enable,attr"`
			URL    string `xml:"url,attr"`
		} `xml:"SafetyDataCollectionURL"`
		SafetyDataCollectionTime struct {
			Text   string `xml:",chardata"`
			Enable string `xml:"enable,attr"`
			URL    string `xml:"url,attr"`
			Time   string `xml:"time,attr"`
		} `xml:"SafetyDataCollectionTime"`
		StatsInfor struct {
			Text          string `xml:",chardata"`
			URL           string `xml:"url,attr"`
			Uploadtime    string `xml:"uploadtime,attr"`
			Interval      string `xml:"interval,attr"`
			Collecttime   string `xml:"collecttime,attr"`
			UpThreshold   string `xml:"upThreshold,attr"`
			DownThreshold string `xml:"downThreshold,attr"`
			Enable        string `xml:"enable,attr"`
		} `xml:"StatsInfor"`
		CheckNetRegularly struct {
			Text         string `xml:",chardata"`
			MaxCount     string `xml:"maxCount,attr"`
			MinBenchmark string `xml:"minBenchmark,attr"`
			IcmpUrl      string `xml:"icmpUrl,attr"`
			HttpUrl      string `xml:"httpUrl,attr"`
			Port         string `xml:"port,attr"`
			Interval     string `xml:"interval,attr"`
			Enable       string `xml:"enable,attr"`
			IcmpTimeout  string `xml:"icmpTimeout,attr"`
			HttpTimeout  string `xml:"httpTimeout,attr"`
			Order        string `xml:"order,attr"`
		} `xml:"CheckNetRegularly"`
		InterfaceInfoUploadURL struct {
			Text   string `xml:",chardata"`
			Enable string `xml:"enable,attr"`
			URL    string `xml:"url,attr"`
		} `xml:"InterfaceInfoUploadURL"`
		WiFiCollectUploadURL struct {
			Text      string `xml:",chardata"`
			Enable    string `xml:"enable,attr"`
			URL       string `xml:"url,attr"`
			Frequency string `xml:"frequency,attr"`
		} `xml:"WiFiCollectUploadURL"`
		CustomerServiceQuestion struct {
			Text   string `xml:",chardata"`
			Enable string `xml:"enable,attr"`
			URL    string `xml:"url,attr"`
		} `xml:"CustomerServiceQuestion"`
		CustomerServiceWelcome struct {
			Text   string `xml:",chardata"`
			Enable string `xml:"enable,attr"`
			URL    string `xml:"url,attr"`
		} `xml:"CustomerServiceWelcome"`
		CustomerServiceAnswer struct {
			Text   string `xml:",chardata"`
			Enable string `xml:"enable,attr"`
			URL    string `xml:"url,attr"`
		} `xml:"CustomerServiceAnswer"`
		CustomerServiceMsgTime struct {
			Text   string `xml:",chardata"`
			Enable string `xml:"enable,attr"`
			Time   string `xml:"time,attr"`
		} `xml:"CustomerServiceMsgTime"`
		UploadClickModule struct {
			Text   string `xml:",chardata"`
			Enable string `xml:"enable,attr"`
			URL    string `xml:"url,attr"`
		} `xml:"UploadClickModule"`
		LocalListenerPort struct {
			Text   string `xml:",chardata"`
			Enable string `xml:"enable,attr"`
			Port   string `xml:"port,attr"`
		} `xml:"LocalListenerPort"`
		Balancereminder struct {
			Text   string `xml:",chardata"`
			Enable string `xml:"enable,attr"`
			URL    string `xml:"url,attr"`
		} `xml:"Balancereminder"`
	} `xml:"funcfg"`
	Emulator struct {
		Text      string `xml:",chardata"`
		Count     string `xml:"count"`
		CheckItem []struct {
			Text            string `xml:",chardata"`
			Name            string `xml:"name,attr"`
			Contains        string `xml:"contains,attr"`
			EqualsCheckItem string `xml:"equals-check-item,attr"`
			CheckFilter     []struct {
				Text    string `xml:",chardata"`
				Model   string `xml:"model,attr"`
				Version string `xml:"version,attr"`
			} `xml:"check-filter"`
		} `xml:"check-item"`
		ExecItem struct {
			Text string `xml:",chardata"`
			Name string `xml:"name,attr"`
		} `xml:"exec-item"`
		SureItem struct {
			Text     string `xml:",chardata"`
			Name     string `xml:"name,attr"`
			Contains string `xml:"contains,attr"`
		} `xml:"sure-item"`
	} `xml:"emulator"`
}

func (cl *Client) GenerateGetTicketXML() ([]byte, error) {
	tr := TicketRequest{
		UserAgent: UserAgentAndroid,
		ClientID:  cl.ClientID.String(),
		LocalTime: time.Now().Format(time.DateTime),
		HostName:  cl.Hostname,
		Ipv4:      cl.UserIP,
		Mac:       cl.MacAddress,
		Ostag:     cl.Hostname,
		Gwip:      cl.AcIP,
	}
	out, err := xml.Marshal(tr)
	if err != nil {
		return nil, err
	}
	return append([]byte("<?xml version=\"1.0\" encoding=\"UTF-8\"?>"), out...), nil
}

func (cl *Client) GenerateStateXML() ([]byte, error) {
	s := &State{
		UserAgent: UserAgentAndroid,
		ClientID:  cl.ClientID.String(),
		LocalTime: time.Now().Format(time.DateTime),
		HostName:  cl.Hostname,
		Ipv4:      cl.UserIP,
		Ticket:    cl.Ticket,
		Mac:       cl.MacAddress,
		Ostag:     cl.Hostname,
	}
	bytes, err := xml.Marshal(s)
	if err != nil {
		return nil, err
	}

	return append([]byte("<?xml version=\"1.0\" encoding=\"UTF-8\"?>"), bytes...), nil
}

func (cl *Client) GenerateLoginXML() ([]byte, error) {
	lr := &LoginRequest{
		UserAgent: UserAgentAndroid,
		ClientID:  cl.ClientID.String(),
		Ticket:    cl.Ticket,
		LocalTime: time.Now().Format(time.DateTime),
		Userid:    cl.Conf.AuthUsername,
		Passwd:    cl.Conf.AuthPassword,
	}

	bytes, err := xml.Marshal(lr)
	if err != nil {
		return nil, err
	}

	return append([]byte("<?xml version=\"1.0\" encoding=\"UTF-8\"?>"), bytes...), nil
}

func (cl *Client) PostXML(url string, data []byte) ([]byte, error) {
	encXML, err := cl.cipher.Encrypt(data)
	if err != nil {
		return nil, err
	}

	req, err := cl.GeneratePostRequest(url, encXML)
	if err != nil {
		return nil, err
	}

	response, err := cl.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	data, err = io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return cl.cipher.Decrypt(data)
}

func (cl *Client) PostXMLWithSpecCtx(url string, data []byte) ([]byte, error) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Second))
	defer cancel()
	encXML, err := cl.cipher.Encrypt(data)
	if err != nil {
		return nil, err
	}

	req, err := cl.GeneratePostRequestWithSpecCtx(ctx, url, encXML)
	if err != nil {
		return nil, err
	}

	response, err := cl.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	data, err = io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return cl.cipher.Decrypt(data)
}
