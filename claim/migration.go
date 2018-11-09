package claim

import (
	"encoding/hex"

	"github.com/btcsuite/btcutil/base58"
	"github.com/lbryio/types/go"
)

const lbrySDHash = "lbry_sd_hash"

func newClaim() *pb.Claim {
	pbClaim := new(pb.Claim)
	stream := new(pb.Stream)
	metadata := new(pb.Metadata)
	source := new(pb.Source)
	pubSig := new(pb.Signature)
	fee := new(pb.Fee)
	metadata.Fee = fee
	stream.Metadata = metadata
	stream.Source = source
	pbClaim.Stream = stream
	pbClaim.PublisherSignature = pubSig

	//Fee version
	feeVersion := pb.Fee__0_0_1
	pbClaim.GetStream().GetMetadata().GetFee().Version = &feeVersion
	//Metadata version
	mdVersion := pb.Metadata__0_1_0
	pbClaim.GetStream().GetMetadata().Version = &mdVersion
	//Source version
	srcVersion := pb.Source__0_0_1
	pbClaim.GetStream().GetSource().Version = &srcVersion
	//Stream version
	streamVersion := pb.Stream__0_0_1
	pbClaim.GetStream().Version = &streamVersion
	//Claim version
	clmVersion := pb.Claim__0_0_1
	pbClaim.Version = &clmVersion
	//Claim type
	clmType := pb.Claim_streamType
	pbClaim.ClaimType = &clmType

	return pbClaim
}

func setMetaData(claim pb.Claim, author string, description string, language pb.Metadata_Language, license string,
	licenseURL *string, title string, thumbnail *string, nsfw bool) {
	claim.GetStream().GetMetadata().Author = &author
	claim.GetStream().GetMetadata().Description = &description
	claim.GetStream().GetMetadata().Language = &language
	claim.GetStream().GetMetadata().License = &license
	claim.GetStream().GetMetadata().Title = &title
	claim.GetStream().GetMetadata().Thumbnail = thumbnail
	claim.GetStream().GetMetadata().Nsfw = &nsfw
	claim.GetStream().GetMetadata().LicenseUrl = licenseURL

}

func migrateV1Claim(vClaim V1Claim) (*pb.Claim, error) {
	pbClaim := newClaim()
	//Not part of json V1
	pbClaim.PublisherSignature = nil
	//Stream
	// -->Universal
	setFee(vClaim.Fee, pbClaim)
	// -->MetaData
	language := pb.Metadata_Language(pb.Metadata_Language_value[vClaim.Language])
	setMetaData(*pbClaim, vClaim.Author, vClaim.Description, language,
		vClaim.License, nil, vClaim.Title, vClaim.Thumbnail, false)
	// -->Source
	pbClaim.GetStream().GetSource().ContentType = &vClaim.ContentType
	sourceType := pb.Source_SourceTypes(pb.Source_SourceTypes_value[lbrySDHash])
	pbClaim.GetStream().GetSource().SourceType = &sourceType
	src, err := hex.DecodeString(vClaim.Sources.LbrySDHash)
	pbClaim.GetStream().GetSource().Source = src

	return pbClaim, err
}

func migrateV2Claim(vClaim V2Claim) (*pb.Claim, error) {
	pbClaim := newClaim()
	//Not part of json V2
	pbClaim.PublisherSignature = nil
	//Stream
	// -->Fee
	setFee(vClaim.Fee, pbClaim)
	// -->MetaData
	language := pb.Metadata_Language(pb.Metadata_Language_value[vClaim.Language])
	setMetaData(*pbClaim, vClaim.Author, vClaim.Description, language,
		vClaim.License, vClaim.LicenseURL, vClaim.Title, vClaim.Thumbnail, vClaim.NSFW)
	// -->Source
	pbClaim.GetStream().GetSource().ContentType = &vClaim.ContentType
	sourceType := pb.Source_SourceTypes(pb.Source_SourceTypes_value[lbrySDHash])
	pbClaim.GetStream().GetSource().SourceType = &sourceType
	src, err := hex.DecodeString(vClaim.Sources.LbrySDHash)
	pbClaim.GetStream().GetSource().Source = src

	return pbClaim, err
}

func migrateV3Claim(vClaim V3Claim) (*pb.Claim, error) {
	pbClaim := newClaim()
	//Not part of json V3
	pbClaim.PublisherSignature = nil
	//Stream
	// -->Fee
	setFee(vClaim.Fee, pbClaim)
	// -->MetaData
	language := pb.Metadata_Language(pb.Metadata_Language_value[vClaim.Language])
	setMetaData(*pbClaim, vClaim.Author, vClaim.Description, language,
		vClaim.License, vClaim.LicenseURL, vClaim.Title, vClaim.Thumbnail, vClaim.NSFW)
	// -->Source
	pbClaim.GetStream().GetSource().ContentType = &vClaim.ContentType
	sourceType := pb.Source_SourceTypes(pb.Source_SourceTypes_value[lbrySDHash])
	pbClaim.GetStream().GetSource().SourceType = &sourceType
	src, err := hex.DecodeString(vClaim.Sources.LbrySDHash)
	pbClaim.GetStream().GetSource().Source = src

	return pbClaim, err
}

func setFee(fee *Fee, pbClaim *pb.Claim) {
	if fee != nil {
		amount := float32(0.0)
		currency := pb.Fee_LBC
		address := ""
		if fee.BTC != nil {
			amount = fee.BTC.Amount
			currency = pb.Fee_BTC
			address = fee.BTC.Address
		} else if fee.LBC != nil {
			amount = fee.LBC.Amount
			currency = pb.Fee_LBC
			address = fee.LBC.Address
		} else if fee.USD != nil {
			amount = fee.USD.Amount
			currency = pb.Fee_USD
			address = fee.USD.Address
		}
		//Fee Settings
		pbClaim.GetStream().GetMetadata().GetFee().Amount = &amount
		pbClaim.GetStream().GetMetadata().GetFee().Currency = &currency
		pbClaim.GetStream().GetMetadata().GetFee().Address = base58.Decode(address)
	}
}
