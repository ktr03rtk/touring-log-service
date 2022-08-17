export type TripTypes = {
  lat: number;
  lng: number;
};

export type PhotoTypes = {
  id: string;
  lat: number;
  lng: number;
};

export type TouringLogTypes = {
  trip: TripTypes[];
  photo: PhotoTypes[];
  center: TripTypes;
};
