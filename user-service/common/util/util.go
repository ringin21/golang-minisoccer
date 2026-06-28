package util

import (
	"os"
	"reflect"
	"strconv"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
)

// BindFromJSON membaca file konfigurasi JSON lokal lalu mengisi data ke dest.
func BindFromJSON(dest any, filename, path string) error {
	v := viper.New()
	v.SetConfigType("json")
	v.AddConfigPath(path)
	v.SetConfigName(filename)

	err := v.ReadInConfig()
	if err != nil {
		return err
	}

	// Unmarshal digunakan untuk memetakan isi config ke struct tujuan.
	err = v.Unmarshal(&dest)
	if err != nil {
		logrus.Errorf("failed to unmarshal : %v", err)
		return err
	}
	return nil
}

// SetEnvFromConsulKV mengubah key-value dari Viper menjadi environment variable.
func SetEnvFromConsulKV(v *viper.Viper) error {
	env := make(map[string]any)

	err := v.Unmarshal(&env)
	if err != nil {
		logrus.Errorf("failed to unmarshal : %v", err)
		return err
	}

	for k, v := range env {
		// Value dari Consul bisa memiliki tipe berbeda, jadi perlu dikonversi ke string sebelum menjadi env.
		var (
			valOf = reflect.ValueOf(v)
			val   string
		)

		switch valOf.Kind() {
		case reflect.String:
			val = valOf.String()
		case reflect.Int:
			val = strconv.Itoa(int(valOf.Int()))
		case reflect.Uint:
			val = strconv.Itoa(int(valOf.Int()))
		case reflect.Float32:
			val = strconv.Itoa(int(valOf.Float()))
		case reflect.Float64:
			val = strconv.Itoa(int(valOf.Float()))
		case reflect.Bool:
			val = strconv.FormatBool(valOf.Bool())
		default:
			panic("unsupported type")
		}

		// Set environment variable berdasarkan key dan value dari konfigurasi.
		err = os.Setenv(k, val)
		if err != nil {
			logrus.Errorf("failed to set env: %v", err)
			return err
		}
	}

	return nil
}

// BindFromConsul membaca konfigurasi JSON dari Consul, mengisi dest, lalu menyet env dari KV Consul.
func BindFromConsul(dest any, endPoint, path string) error {
	v := viper.New()
	v.SetConfigType("json")

	// Remote provider digunakan agar Viper bisa membaca konfigurasi dari Consul.
	err := v.AddRemoteProvider("consul", endPoint, path)
	if err != nil {
		logrus.Errorf("failed to add remote provider: %v", err)
		return err
	}

	err = v.ReadRemoteConfig()
	if err != nil {
		logrus.Errorf("failed to read remote config: %v", err)
		return err
	}

	// Unmarshal digunakan untuk memetakan config Consul ke struct tujuan.
	err = v.Unmarshal(&dest)
	if err != nil {
		logrus.Errorf("failed to unmarshal : %v", err)
	}

	// Setelah config dibaca, key-value Consul juga disimpan ke environment variable.
	err = SetEnvFromConsulKV(v)
	if err != nil {
		logrus.Errorf("failed to set env from consul kv: %v", err)
		return err
	}

	return nil
}

/*
Kegunaan file:
File ini dibuat untuk menyimpan helper konfigurasi yang dipakai bersama oleh aplikasi.
Dengan file ini, aplikasi bisa membaca konfigurasi dari file JSON lokal atau dari Consul,
lalu mengubah konfigurasi tertentu menjadi environment variable yang bisa dipakai bagian lain.
*/
