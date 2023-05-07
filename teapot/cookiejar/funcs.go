package cookiejar

func GetErrorChannel(jar *cookieContainer) ErrorChannel {
	return jar.errchan
}
