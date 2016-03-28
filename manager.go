package main

import "sync"

/*
 Simple interface that allows us to switch out both implementations of the Manager
*/
type ConfigManager interface {
	Set(*Config)
	Get() *Config
	Close()
}

/*
 This struct manages the configuration instance by
 preforming locking around access to the Config struct.
*/
type MutexConfigManager struct {
	conf  *Config
	mutex *sync.Mutex
}

func NewMutexConfigManager(conf *Config) *MutexConfigManager {
	return &MutexConfigManager{conf, &sync.Mutex{}}
}

func (self *MutexConfigManager) Set(conf *Config) {
	self.mutex.Lock()
	self.conf = conf
	self.mutex.Unlock()
}

func (self *MutexConfigManager) Get() *Config {
	self.mutex.Lock()
	temp := self.conf
	self.mutex.Unlock()
	return temp
}

func (self *MutexConfigManager) Close() {
	//Do Nothing
}

/*
 This struct manages the configuration instance by feeding a
 pointer through a channel whenever the user calls Get()
*/
type ChannelConfigManager struct {
	conf *Config
	get  chan *Config
	set  chan *Config
	done chan bool
}

func NewChannelConfigManager(conf *Config) *ChannelConfigManager {
	parser := &ChannelConfigManager{conf, make(chan *Config), make(chan *Config), make(chan bool)}
	parser.Start()
	return parser
}

func (self *ChannelConfigManager) Start() {
	go func() {
		defer func() {
			close(self.get)
			close(self.set)
			close(self.done)
		}()
		for {
			select {
			case self.get <- self.conf:
			case value := <-self.set:
				self.conf = value
			case <-self.done:
				return
			}
		}
	}()
}

func (self *ChannelConfigManager) Close() {
	self.done <- true
}

func (self *ChannelConfigManager) Set(conf *Config) {
	self.set <- conf
}

func (self *ChannelConfigManager) Get() *Config {
	return <-self.get
}
